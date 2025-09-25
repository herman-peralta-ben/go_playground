package main

/// go run main.go

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	// This will be called when main ends
	defer func() {
		fmt.Println("finished main")
	}()

	// [chan] declare channels to return data from goroutine
	cs1 := make(chan string)
	cs2 := make(chan string)

	// [go] executes function as GoRoutine
	// Fire anf forget because is not using a channel to return data
	go func() {
		fmt.Println("Fire and forget goroutine")
	}()

	// executing as goroutine, returns data using [cs2], if someone is listening to it
	go func() {
		time.Sleep(500 * time.Millisecond)
		cs2 <- "two"
	}()

	// Will execute on goroutine and return data using [cs1]
	go sleep1(cs1)

	fmt.Println("Waiting for [select] first channel message with select, which is [cs2]")
	select {
	case msg1 := <-cs1:
		fmt.Println("cs1 received", msg1)
	case msg2 := <-cs2:
		fmt.Println("cs2 received", msg2)
	}
	fmt.Println("After select")

	// Waiting for an specifc channel
	cs1Msg := <-cs1
	fmt.Println("cs1 received", cs1Msg)

	// Simulate streams
	// non blocking
	var wg sync.WaitGroup
	wg.Add(1)
	chNonBlockingStream := make(chan int)
	go buildPeriodicStream(chNonBlockingStream, 100, 50)
	consumeNonBlocking(chNonBlockingStream, &wg)

	// blocking
	ch_str1 := make(chan int)
	go buildPeriodicStream(ch_str1, 300, 5)
	consumeBlocking(ch_str1)

	// Wait for [chNonBlockingStream] to finish
	wg.Wait()
	fmt.Println("=>NonBlocking stream done, WG")
}

func sleep1(ch chan string) {
	time.Sleep(1 * time.Second)
	ch <- "one" // Needs to use a channel to return data
	close(ch)
}

func buildPeriodicStream(ch chan int, millis int, count int) {
	defer close(ch) // close the chanel even if the goroutine has an error
	for i := 0; i < count; i++ {
		ch <- i
		time.Sleep(time.Duration(millis) * time.Millisecond)
	}
}

func consumeBlocking(ch chan int) {
	fmt.Println("=>Blocking stream launched")
	for val := range ch {
		fmt.Println("consumeBlocking:", val)
	}
	fmt.Println("=>Blocking stream done")
}

func consumeNonBlocking(ch chan int, wg *sync.WaitGroup) {
	fmt.Println("=>NonBlocking stream launched")
	go func() {
		defer wg.Done()
		for val := range ch {
			fmt.Println("consumeNonBlocking", val)
		}
	}()
	fmt.Println("=>NonBlocking stream done")
}
