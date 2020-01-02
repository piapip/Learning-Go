Closing a channel indicates that no more values will be sent on it. This can be useful to communicate completion to the channel's receivers

each context.Context should represent for only one data query, not for an entire database object

Page 27:
i++
If your context is a goroutine that doesn't expose i to other goroutines, then this code is atomic.

Page 158: 
Golang doesn't have throw/catch mechanic so you need to build a struct if you want to keep track of possible errors:
type Result struct {
    Error    error
    Response interface{} //change this one
}

Page 162:
Pipeline is nothing too amazing tbh. It's just a bunch of functions return the same type of value, so the function can be "recycled" or "called" by another function 
as variable including itself. One easy example: fibonacci without recursive. 2 types: stream process and batch process

Batch process: 
The function returns a batch, an array of values. Increase scalability, though it will eat up more memory initially, but it makes your works a lot more easier 
later on. (Because everything the function is triggered, it will load the entire *batch* that you just pass into them)

Stream process:
The function returns a single individual, it will greatly reduce the memory required (well, most of the time). But it will render the scalability. 
Not recommended tbh.

Page 169-170 is important.

interface{} is not any type. It will be downcast to the first type that it's assigned to.
A pointer type can access the methods of its associated value type, but not vice versa
(d Dog)doSomething() int {} => (*aDog).doSomething(); OK
(d *Dog)doSomething() int {} => aDog.doSomething(); FAIL
You can use doSomething(a ...interface{}) (1) but you should not use doSomething(a []interface{}) (2) because 
the code can run properly in (1) (if you write doSomething(a int, b string, c bool)) but
the code will return a compile warning in (2) (if you write doSomething( a []int)) Gotta convert everything of that a to interface first
    vals := make([]interface{}, len(names))
	for i, v := range names {
	    vals[i] = v
    }

Why is it not recommended to use {}interface?

In generators, it's fine to use generic indicator (the chan interface{}). Though, you can double the speed with specific indicator (like chan int or chan string) 
but the difference is negligible.

Fan-out is a term to describe the process of starting multiple goroutines to handle input from the pipeline, and fan-in is a term to describe the process of
combining multiple results into one channel. Use fanning if both of this requirements are met:
One of the stages in your pipeline doesn’t rely on values that the stage had calculated before.
And it takes a long time to run.

Fan-in is like a double-edge knife. Could speed up the process but still might cause unpleasant experience because it takes over all cores of the CPU which means
it will lower other applications quality which is not a good thing. (Try playing hots on a toaster and you'll see). Nothing is free, it's just a trade-off. (
Workflow:
Chan1:              Chan2:                  Time consumed of non-fanning: 2(***************) + some stupid (doesn't matter material)
**   1              **      1               ( [**1] + [**2] (chan2 [**1] might be close to finish for now) [***************] +  [**2] + [***************] + the rest)
**   2              **      2               Time consumed of fanning    :  (***************) + some stupid (doesn't matter material, a lil bit more than the non-fanning
**************      **************                                          but unsignificant)
**                  **
**                  **
**                  **
**                  **
**                  **
**                  **
)

So to prevent goroutine leak, you do this in every goroutine: 
instead of:
    for val := range myChan (might leak goroutine) 
TO
    go func() {             (no more leaking)
        for{
            select {
            case <-done:
                return
            case maybeVal, ok := <- myChan: 
                ... 
            }
        }
    }

^ this could get worse if ... is a nested loop because you'll need to "case <-done" again and again. This is like callback hell in js. So we could combine 
or-channel technique to handle this mess. (Like how js has Promise for the callback hell).
Solution:
go func() {
    defer close(valStream)
    for {
        select {
        case <- done: return
        case v, ok := <- valStream:
            if ok == false {
                return
            }
            select {            (cleaner)
            case valStream <- v:
            case <-done:                
            }
        }
    }
}()

and then you just need to do for val := range or_done(done, thatChannelYouWantToTraverseWithoutLeaking)

tee channel... it splits a channel to multiple channels. (the origin doesn't change though)

a lil difference between fan-in and or-channel
fan-in = all the elements that other channels contain

or-channel = all the SIGNAL that other channels contain and SHARE THAT SIGNAL to every channel that is link to the or-channel.

so when channel1 is closed, every channel that is linked to the or-channel as well as the or-channel will be closed, but the fan-in channel won't be.

Can't see much use of bridge channel. If you want to a stream around,
EX: chanStream <- stream. Remember to close stream first, otherwise they will just stand there and block off the flow.

Read 196 for somthing about bufio and queueing.

Context explanation: 201

All context function will always receive a parent Context and split out another Context.
context.Deadline: 215
context usage: 218 + 219, sounds like making loose cohesion.

iota: https://yourbasic.org/golang/iota/

StackTrace: 224

Heartbeats that occur on a time interval are useful for concurrent code that might be waiting for something else to happen for it to process a unit of work. 
Because you don’t know when that work might come in, your goroutine might be sitting around for a while waiting for something to happen. 
A heartbeat is a way to signal to its listeners that everything is well, and that the silence is expected.
250+251 for further details.

In fact, there’s a name for a section of your program that needs exclusive access to a shared resource. This is called a critical section. In this example,
we have three critical sections:
Our goroutine, which is incrementing the data variables.
Our if statement, which checks whether the value of data is 0.
Our fmt.Printf statement, which retrieves the value of data for output.

In Go, goroutines are tasks. If a goroutine spawns other goroutines then it's spawning tasks, if a goroutine spawns a non-goroutine function then it's just 
minding its own business.
Everything after a goroutine is called is the continuation.
Go’s work-stealing algorithm enqueues and steals continuations. 
When a thread of execution reaches an unrealized join point, the thread must pause execution and go fishing for a task to steal.