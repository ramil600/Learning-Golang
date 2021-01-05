package main

import (
	"crypto/md5"
	"fmt"
	"hash/crc32"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

type job func(in, out chan interface{})

var done chan interface{}

var (
	dataSignerOverheat uint32 = 0
	DataSignerSalt            = ""
)
var DataSignerCrc32 = func(data string) string {
	data += DataSignerSalt
	crcH := crc32.ChecksumIEEE([]byte(data))
	dataHash := strconv.FormatUint(uint64(crcH), 10)
	time.Sleep(time.Second)
	return dataHash
}
var DataSignerMd5 = func(data string) string {
	OverheatLock()
	defer OverheatUnlock()
	data += DataSignerSalt
	dataHash := fmt.Sprintf("%x", md5.Sum([]byte(data)))
	time.Sleep(10 * time.Millisecond)
	return dataHash
}

var OverheatLock = func() {
	for {
		if swapped := atomic.CompareAndSwapUint32(&dataSignerOverheat, 0, 1); !swapped {
			fmt.Println("OverheatLock happend")
			time.Sleep(time.Second)
		} else {
			break
		}
	}
}

var OverheatUnlock = func() {
	for {
		if swapped := atomic.CompareAndSwapUint32(&dataSignerOverheat, 1, 0); !swapped {
			fmt.Println("OverheatUnlock happend")
			time.Sleep(time.Second)
		} else {
			break
		}
	}
}

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

var hashonce = func(in string, num int) string {

	concated := strconv.Itoa(num) + in

	return DataSignerCrc32(concated)

}

func SingleHash(in chan interface{}, out chan interface{}) {
	var a, b string
	var wg sync.WaitGroup
	go func() {
		defer close(out)
		for i := range in {
			j := strconv.Itoa(i.(int))
			wg.Add(2)
			go func() {
				defer wg.Done()
				a = DataSignerCrc32(j)
			}()
			go func() {
				defer wg.Done()
				b = DataSignerCrc32(DataSignerMd5(j))
			}()

			wg.Wait()
			c := a + "~" + b

			go func() {
				select {
				case out <- c:
				case <-done:
				}
			}()

		}
	}()

	for e := 0; e < 5; e++ {
		fmt.Println(<-out)
	}

}

func MultiHash(in chan interface{}, out chan interface{}) {

	finders := make([]chan interface{}, 6)
	var wg sync.WaitGroup
	var b string
	for data := range in {
		wg.Add(6)
		for j := 0; j < 6; j++ {
			go func(j int) {
				defer wg.Done()
				finders[j] <- hashonce(data.(string), j)

			}(j)
		}

		go func() {
			wg.Wait()
			for i := 0; i < 6; i++ {
				a := <-finders[i]
				b += a.(string)

			}
			select {
			case out <- b:
			case <-done:
				return
			}

		}()

	}

}

//first function generates stream of ints, it is feeding the pipeline
func generator(in chan interface{}, out chan interface{}) {
	ints := []int{0, 1, 2, 3, 4, 5}
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
	ExecutePipeline(generator, SingleHash)
	fmt.Println(time.Since(start), "elapsed since the start of the pipeline.")

	//fmt.Println(DataSignerCrc32("0"), "~", DataSignerCrc32(DataSignerMd5("0")))
}
