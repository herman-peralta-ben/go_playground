package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("========= synchronousGoroutine =========")
	synchronousGoroutine()
	fmt.Println("========= asynchronousGoroutine =========")
	asynchronousGoroutine()
}

func synchronousGoroutine() {
	// capacity 1
	blockingChannel := make(chan int)

	// sender
	go func() {
		fmt.Println("Sender: before sending")
		blockingChannel <- 42
		// This will be printed after reading the message
		fmt.Println("Sender: after sending")
	}()

	fmt.Println("synchronousGoroutine: before receiving")
	time.Sleep(1 * time.Second)
	msg := <-blockingChannel
	fmt.Printf("synchronousGoroutine: after receiving, got %d\n", msg)
}

func asynchronousGoroutine() {
	// capacity 2
	blockingChannel := make(chan int, 2)

	// sender
	go func() {
		fmt.Println("Sender: before sending")
		blockingChannel <- 42
		// This will be printed even if the message is not consumed
		fmt.Println("Sender: after sending, even if message is not consumed from channel")
	}()

	fmt.Println("asynchronousGoroutine: before sleeping")
	time.Sleep(1 * time.Second)
	fmt.Println("asynchronousGoroutine: before sleeping")
}
