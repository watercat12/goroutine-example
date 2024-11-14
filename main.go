package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	// unbuffered()
	// buffered()
	// waitGroupExample()
	// noCloseChannel()
}

func noCloseChannel() {
	var wg sync.WaitGroup
	errChan := make(chan error, 4)
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			time.Sleep(time.Duration(i) * time.Second)
			errChan <- fmt.Errorf("error %d", i)
		}(i)
	}

	wg.Wait()
	for i := 0; i < 4; i++ {
		if err := <-errChan; err != nil {
			fmt.Println(err)
		}
	}
	// this code will cause deadlock
	// for err := range errChan {
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	}
	// }
}

func waitGroupExample() {
	var wg sync.WaitGroup
	errChan := make(chan error)
	// we can using wg.Add(4) here, then remove the wg.Add(1) in the loop
	// wg.Add(4) to add 4 goroutines to the wait group
	wg.Add(3) // try using wg.Add(3) to see the difference
	for i := 0; i < 4; i++ {
		// wg.Add(1)
		go func(i int) {
			defer wg.Done() // when wg.Done() is called, the counter of wg is decremented by 1
			time.Sleep(time.Duration(i) * time.Second)
			if i < 3 {
				errChan <- fmt.Errorf("error %d", i)
			}
		}(i)
	}

	go func() {
		wg.Wait()
		fmt.Println("all done")
		close(errChan)
	}()

	for err := range errChan {
		if err != nil {
			fmt.Println(err)
		}
	}
	fmt.Println("done")
}

func buffered() {
	var wg sync.WaitGroup
	errChan := make(chan error, 2)

	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			time.Sleep(time.Duration(i) * time.Second)
			if i < 3 {
				errChan <- fmt.Errorf("error %d", i)
			}
		}(i)
	}
	go func() {
		wg.Wait()
		close(errChan)
	}()

	for err := range errChan {
		if err != nil {
			fmt.Println(err)
		}
	}
}

func unbuffered() {
	var wg sync.WaitGroup
	errChan := make(chan error)

	for i := 0; i < 4; i++ {
		wg.Add(1)
		fmt.Println("wait group add ", i)
		go func(i int) {
			defer func() {
				fmt.Println("goroutine done ", i)
				wg.Done()
			}()
			time.Sleep(time.Duration(i) * time.Second)
			if i < 3 {
				errChan <- fmt.Errorf("error %d", i)
				fmt.Printf("Add %d to channel\n", i)
			}
		}(i)
	}

	go func() {
		fmt.Println("wg.Wait()")
		wg.Wait()
		fmt.Println("done wg.Wait()")
		close(errChan)
	}()

	for err := range errChan {
		if err != nil {
			fmt.Println("--------------------------------")
			fmt.Printf("|main receive error %v|\n", err)
			fmt.Println("--------------------------------")
		}
	}
	fmt.Println("main done")
}

// The main differences are:

// 1. The channel `errChan` is now created as an unbuffered channel: `errChan := make(chan error)`.
// 2. A new goroutine is created to call `wg.Wait()` and then close the channel after all goroutines have completed.
// This is necessary because with an unbuffered channel,
// the send operation `errChan <- fmt.Errorf("error %d", i)`
// will block until there's a corresponding receive operation on the other end of the channel.

// In the original code with a buffered channel,
// the send operations could complete without blocking since the buffer had a capacity of 2.
//  With an unbuffered channel,
//  the send operations would block until the `for err := range errChan` loop starts receiving values from the channel.

// By creating a separate goroutine to call `wg.Wait()` and close the channel,
// we ensure that the send operations in the goroutines can complete,
// and the `for err := range errChan` loop can read all the values from the channel before terminating.
