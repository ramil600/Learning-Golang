//Cond  structure with sync package.
// We need to add 10 elements, we want max 2 elements in the queue, so we call a removeq goroutine
//that will remove elements from the queue one length becomes 2.
//we use L lock that belongs to Cond c to Lock and Unlock in critical sections of main and removeq
//goroutine that removes one element at a time
//c.Wait() waits on goroutine while length is equal to 2

package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {

	queue := make([]interface{}, 0, 10)
	
        //We create NewCond with sync.Mutex struct to control the length of the queue
  	c := sync.NewCond(&sync.Mutex{})
	
	removeq := func(){
		time.Sleep(time.Second)
		c.L.Lock()
		queue = queue[1:]
		fmt.Println("Removing one from queue")
		c.L.Unlock()
		c.Signal()
	}
	
	for i:=0; i <10; i++{
		c.L.Lock()
		
                //if we reach code where the length is 2, element needs to be removed before we proceed
	        //we then wait in the loop suspended, L lock is released, this avoids loading on CPU. When removeq 
                //signals that element removed, we can then restore L lock, proceed with adding onother element
		for len(queue) == 2 {
			c.Wait()
		}
		
		queue = append(queue, struct{}{})
		fmt.Println("Adding one to queue")
		go removeq()

		c.L.Unlock()
	}
}
