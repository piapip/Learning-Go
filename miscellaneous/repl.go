package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"contrib.go.opencensus.io/exporter/prometheus"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

// So this is metrics :) srsl you can name it w/e you want. I'm just blindly following the guide.
var (
	// MLatencyMs latency milliseconds
	MLatencyMs = stats.Float64("relp/latency", "The latency in milliseconds per REPL loop", "ms")

	// MLineLengths - Counts/groups the lengths of lines read in
	MLineLengths = stats.Int64("repl/line_lengths", "The distribution of line lengths", "By")
)

var (
	//KeyMethod record what method is being invoked, it will say that "relp" is calling our data
	KeyMethod, _ = tag.NewKey("method")
	//KeyStatus record the status after each calls.
	KeyStatus, _ = tag.NewKey("status")
	//KeyError self-explanatory
	KeyError, _ = tag.NewKey("error")
)

//how our metrics will be organized
var (
	//LatencyView will display latency
	LatencyView = &view.View{
		Name:        "demo/latency",
		Measure:     MLatencyMs,
		Description: "The distribution of the latencies",

		// Latency in buckets:
		// [>=0ms, >=25ms, >=50ms, >=75ms, >=100ms, >=200ms, >=400ms, >=600ms, >=800ms, >=1s, >=2s, >=4s, >=6s]
		Aggregation: view.Distribution(0, 25, 50, 75, 100, 200, 400, 600, 800, 1000, 2000, 4000, 6000),
		TagKeys:     []tag.Key{KeyMethod},
	}

	LineCountView = &view.View{
		Name:        "demo/lines_in",
		Measure:     MLineLengths,
		Description: "The number of lines from standard input",
		Aggregation: view.Count(),
	}

	LineLengthView = &view.View{
		Name:        "demo/line_lengths",
		Description: "Groups the lengths of keys in buckets",
		Measure:     MLineLengths,
		// Lengths: [>=0B, >=5B, >=10B, >=15B, >=20B, >=40B, >=60B, >=80, >=100B, >=200B, >=400, >=600, >=800, >=1000]
		Aggregation: view.Distribution(0, 5, 10, 15, 20, 40, 60, 80, 100, 200, 400, 600, 800, 1000),
	}
)

func main() {
	// Register the views, it is imperative that this step exists
	// lest recorded metrics will be dropped and never exported.
	if err := view.Register(LatencyView, LineCountView, LineLengthView); err != nil {
		log.Fatalf("Failed to register views: %v", err)
	}

	// Create the Prometheus exporter.
	pe, err := prometheus.NewExporter(prometheus.Options{
		Namespace: "ocmetricstutorial",
	})
	if err != nil {
		log.Fatalf("Failed to create the Prometheus stats exporter: %v", err)
	}

	// Now finally run the Prometheus exporter as a scrape endpoint.
	// We'll run the server on port 8888.
	go func() {
		mux := http.NewServeMux()
		mux.Handle("/metrics", pe)
		if err := http.ListenAndServe(":8888", mux); err != nil {
			log.Fatalf("Failed to run Prometheus scrape endpoint: %v", err)
		}
	}()

	//In a REPL:
	//  1.Read input
	//  2. Process input
	br := bufio.NewReader(os.Stdin)

	// Register the views
	if err := view.Register(LatencyView, LineCountView, LineLengthView); err != nil {
		log.Fatalf("Failed to register views: %v", err)
	}

	// repl is the read, evaluate, print, loop
	for {
		if err := readEvaluateProcess(br); err != nil {
			if err == io.EOF {
				return
			}
			log.Fatal(err)
		}
	}
}

func readEvaluateProcess(br *bufio.Reader) (terr error) {
	//ctx record our metrics
	ctx, err := tag.New(context.Background(), tag.Insert(KeyMethod, "relp"), tag.Insert(KeyStatus, "OK"))
	if err != nil {
		return err
	}

	//just in case it get fucked
	defer func() {
		if terr != nil {
			ctx, _ = tag.New(ctx, tag.Upsert(KeyStatus, "ERROR"), tag.Upsert(KeyError, terr.Error()))
		}

		// stats.Record(ctx, MLatencyMs.M(sinceInMilliseconds(startTime)))
		stats.Record(ctx, MLatencyMs.M(sinceInMilliseconds(time.Now())))
	}()

	fmt.Printf("> ")
	line, _, err := br.ReadLine()
	if err != nil {
		return err
	}

	//make ctx available to record
	out, err := processLine(ctx, line)
	if err != nil {
		return err
	}
	fmt.Printf("< %s\n\n", out)
	return nil
}

//recording metrics here
func processLine(ctx context.Context, in []byte) (out []byte, err error) {
	startTime := time.Now()
	defer func() {
		//LUL settings up matrics is complicated
		stats.Record(ctx, MLatencyMs.M(sinceInMilliseconds(startTime)), MLineLengths.M(int64(len(in))))
	}()

	return bytes.ToUpper(in), nil
}

func sinceInMilliseconds(startTime time.Time) float64 {
	return float64(time.Since(startTime).Nanoseconds()) / 1e6
}
