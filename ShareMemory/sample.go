//So what's the point of using pointer in this piece of code, it still works fine without using pointer
//Buffered channel vs unbuffered channel
//Buffered: Unlimited queue items, needs to be run in a different Goroutine
//Unbuffered: Limited queue items, can run anywhere
package main

import (
	"log"
	"net/http"
	"time"
)

const (
	numPollers     = 2               // number of Poller goroutines to launch
	pollInterval   = 5 * time.Second // how often to poll each URL
	statusInterval = 1 * time.Second // how often to log status to stdout
	errTimeout     = 1 * time.Second // back-off timeout on error
)

var urls = []string{
	"http://www.google.com/",
	"http://golang.org/",
	"http://blog.golang.org/",
}

// State represents the last-known state of a URL
type State struct {
	url    string
	status string
}

//Use log instead of fmt when the programme involves using Goroutine
func logState(s map[string]string) {
	log.Println("Current state: ")
	for k, v := range s {
		log.Printf(" %s %s", k, v)
	}
}

// StateMonitor maintains a map that stores the state of the URLs being
// polled, and prints the current state every updateInterval nanoseconds.
// It returns a chan State to which resource state should be sent.
func StateMonitor(updateInterval time.Duration) chan<- State {
	updates := make(chan State)
	urlStatus := make(map[string]string) //urlStatus[url] = status (200 OK or 400 FAIL)
	ticket := time.NewTicker(updateInterval)
	go func() {
		for {
			select {
			case <-ticket.C: //print urlStatus every updateInterval
				logState(urlStatus)
			//initiate the urlStatus, only trigger once at the beginnning of the code, here's an example for s.status:{http://www.google.com/ 200 OK}
			case s := <-updates:
				urlStatus[s.url] = s.status
			}
		}
	}()
	return updates
}

// Resource represents an HTTP URL to be polled by this program.
type Resource struct {
	url      string
	errCount int
}

// Poll executes an HTTP HEAD request for url
// and returns the HTTP status string or an error string.
func (r *Resource) Poll() string {
	//Neglect this crap for now
	resp, err := http.Head(r.url)
	if err != nil {
		log.Println("Error", r.url, err)
		r.errCount++
		return err.Error()
	}
	r.errCount = 0
	return resp.Status
}

// Sleep sleeps for an appropriate interval (dependent on error state), which means that it only trigger if you give this programme shit url
// before sending the Resource to done.
// It runs in a Goroutine, and then it will return the value to that Goroutine.
func (r *Resource) Sleep(done chan<- *Resource) {
	time.Sleep(pollInterval + errTimeout*time.Duration(r.errCount))
	done <- r
}

//Poller do this
func Poller(in <-chan *Resource, out chan<- *Resource, status chan<- State) {
	//so here is an example for r: 2019/10/30 17:17:45 {http://www.google.com/ 0}
	for r := range in {
		// log.Println(*r)
		// s will update the status of the url
		s := r.Poll()
		status <- State{r.url, s}
		out <- r
	}
}

func main() {
	k := 0
	// Create our input and output channels.
	pending, complete := make(chan *Resource), make(chan *Resource)

	// Launch the StateMonitor.
	// a channel that hold the information about those urls
	status := StateMonitor(statusInterval)

	// Launch some Poller goroutines.
	for i := 0; i < numPollers; i++ {
		go func() {
			Poller(pending, complete, status)
			k++
			log.Println("Under Poller...", k)
		}()
	}

	// Send some Resources to the pending queue.
	go func() {
		for _, url := range urls {
			pending <- &Resource{url: url}
			k++
			log.Println("Under sending &Resource...", k)
		}
	}()

	for r := range complete {
		//also got in here once
		//waiting for error to snitch then sleep if you give it shit url then it will get in here a couple times
		go r.Sleep(pending)
		k++
		log.Println("Under sleeping...", k)
	}

}
