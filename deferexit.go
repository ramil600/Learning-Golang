// You want to clock slow operation on entry and on exit. It is implemented here by defering the function that is returned by anonymous 
// function below. Please pay attention to two sets of paranthesis, first calls anonymous function second one is for return function call,
//  that is deferred. It will be eventually called on exit, meanwhile function itself will be called at entry.
package main

import (
	"fmt"
	"time"
)

func slowoperation() {
  	defer func() func() {
		start := time.Now()
    	// Executed when function slowoperation starts
		fmt.Println("Started at",start)
		return func() {
      	// Return function will be deferred until slowoperation is done.
			fmt.Println("Finished after", time.Since(start))
		}
  // Pay attention to 2 sets of paranthesis. Return function call will be deferred, if we put only one set, return function will never be called.	
  }()()
  // Imitating some slow job
  time.Sleep(5 * time.Second)

}

func main() {
	slowoperation()
	fmt.Println("Bye bye")
}
