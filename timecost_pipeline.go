package main

import (
	"fmt"
	"time"
)

func main() {
	
	start := time.Now()

	//this is the stage of the pipeline closes  out channel when exits, triggers the functions down the
	//pipeline to stop reading from out channel, has done channel that triggers to exit if main exits
	multiply := func(done <-chan interface{}, integers <-chan int, mult int) <-chan int {
		out := make(chan int)
		go func() {
			defer close(out)
			for i := range integers {
				select {
				case <-done:
					return
				case out <- i * mult:
					time.Sleep(3 * time.Millisecond)
				}
			}
		}()
		return out
	}
	
  	//similar stage of the pipeline as multiply, exept costs 4ms to run
	add := func(done <-chan interface{}, integers <-chan int, add int) <-chan int {
		out := make(chan int)
		go func() {
			defer close(out)
			for i := range integers {
				select {
				case <-done:
					return
				case out <- i + add:
					time.Sleep(4 * time.Millisecond)
				}
			}
		}()
		return out
  	}
  	
	//generate the feed to our pipeline
	generator := func(integers ...int) <-chan int {
		out := make(chan int)
		go func() {
			defer close(out)
			for _, i := range integers {
				out <- i
			}
		}()
		return out
	}

  	inputs := generator(1, 2, 3, 4)

  	//avoid leakage of goroutines, close the channel when main func exits
	done := make(chan interface{})
	defer close(done)

  	//pipeline of three costly operations: add 4ms, multiply 3ms
	outputs := multiply(done, add(done, multiply(done, inputs, 2), 1), 3)

	for j := range outputs {
		fmt.Println(j)
	}
	//calculate time since the start of the pipeline
  	dur := time.Since(start)
  
  	//total time cost of our pipeline is only the total cost of 16ms, more efficient that sequential cost of 40ms
	fmt.Println("Total cost of operation is: ", dur.Milliseconds(), "milliseconds.")
}
