package main

import (
	"fmt"
	"time"
)

type job func(in, out chan interface{})

var done chan interface{}

//takes in variadic parameters of job type functions, implements pipeline of the jobs
func ExecutePipeline(hashjobs ...job) {
	outs := make([]chan interface{}, len(hashjobs))

	//I implemented first function outside of the for loop because it does not really takes input channel
	var in chan interface{}
	outs[0] = make(chan interface{})
	hashjobs[0](in, outs[0])

	//For every job in the pipeline, except first, the output channel becomes the input channel for the next job
	for i, hashjob := range hashjobs {
		if i != 0 {
			outs[i] = make(chan interface{})
			hashjob(outs[i-1], outs[i])
		}
	}
}

//first function generates stream of ints, it is feeding the pipeline
func generator(in chan interface{}, out chan interface{}) {
	ints := []int{1, 2, 3, 4, 5}
	go func() {
		defer close(out)
		defer fmt.Println("Generator Exited")

		for _, i := range ints {
			select {
			case out <- i:
			case <-done:
			}
		}
	}()
}

func adder(in chan interface{}, out chan interface{}) {
	go func() {
		defer close(out)
		defer fmt.Println("Adder Exited")

		for i := range in {
			time.Sleep(time.Second)
			select {
			case out <- i.(int) + 1:
			case <-done:
			}
		}
	}()
}

func multiplier(in chan interface{}, out chan interface{}) {
	go func() {
		defer close(out)
		defer fmt.Println("Multiplier Exited")

		for i := range in {
			time.Sleep(2 * time.Second)
			select {
			case out <- i.(int) * 2:
			case <-done:
			}
		}
	}()
	//spits out the integers from the output channel of the last function in the pipeline
	for j := range out {
		fmt.Println(j)
	}
}

func main() {

	//done is implemented as a measure to avoid go routines leak in case the main program exits
	done = make(chan interface{})
	defer close(done)

	/*  You can alternatively pass hashjobs as an array of jobs to ExecutePipeline
	hashjobs := []job{
		job(generator),
		job(adder),
		job(multiplier),
	}
	*/
	start := time.Now()
	ExecutePipeline(generator, adder, multiplier)
	fmt.Println(time.Since(start), "elapsed since the start of the pipeline.")
	fmt.Println("Every adder needs 1sec, whereas multiplier 2 secs; total 5 integers are done in 11sec")
}
