// This program demonstrates the "OR" func which takes several channels(works) that takes different time to complete..
// and returns channel(unblocks) from the fastest task.
// We simulate task with help of sig func, that will unblock after certain delay


package main

import (
	"fmt"
	"time"
)

func main() {
	
	var or func(channels ...<-chan interface{}) <-chan interface{}

	or = func(channels ...<-chan interface{}) <-chan interface{} {
		done := make(chan interface{})

		switch len(channels) {
		case 0:
			return nil
		case 1:
			return channels[0]

		}

		go func() {
			defer close(done)

			switch len(channels) {

			case 2:
				select {
				case <-channels[0]:
				case <-channels[1]:
				}
				fmt.Println("Case 2")

			default:
				fmt.Println("Case default")
				select {
				case <-channels[0]:
					fmt.Println("Read from channel 0")
				case <-channels[1]:
					fmt.Println("Read from channel 1")
				case <-or(append(channels[2:], done)...):
					fmt.Println("Recursive or function")

				}

			}

		}()

		return done

	}

	sig := func(delay time.Duration) <-chan interface{} {

		ordone := make(chan interface{})

		go func() {
			defer close(ordone)
			time.Sleep(delay)
			fmt.Println("Sleep for", delay)

		}()

		return ordone
	}

	<-or(sig(time.Second),
		sig(time.Minute),
		sig(time.Hour))
}
