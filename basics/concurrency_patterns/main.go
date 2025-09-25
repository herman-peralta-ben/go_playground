package main

import (
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"
)

// Part1 https://www.youtube.com/watch?v=qyM8Pi1KiiM
// Part2 https://www.youtube.com/watch?v=wELNUHb3kuA
func main() {
	const exampleCount = 10
	if len(os.Args) < 2 {
		fmt.Printf("Provide an exampleId between 1 and %d\n", exampleCount)
		os.Exit(1)
	}

	arg := os.Args[1]
	exampleId, err := strconv.Atoi(arg)
	if err != nil || exampleId < 1 || exampleId > exampleCount {
		fmt.Printf("Invalid arg '%v'\n", err)
		os.Exit(1)
	}

	switch exampleId {
	case 1:
		fmt.Println("========= testingGoForkJoinModel =========")
		testingGoForkJoinModel()
	case 2:
		fmt.Println("========= blockingChannel =========")
		blockingChannel()
	case 3:
		fmt.Println("========= nonBlockingChannel =========")
		nonBlockingChannel()
	case 4:
		fmt.Println("========= blockingChannelWithGoRoutine =========")
		blockingChannelWithGoRoutine()
	case 5:
		fmt.Println("========= forSelectPattern =========")
		forSelectPattern()
	case 6:
		fmt.Println("========= asyncCharGoroutine2 =========")
		asyncCharGoroutine2()
		// TODO 1. for-select pattern
		// TODO 2. done channel pattern
		// TODO 3. pipelines
	case 7:
		fmt.Println("========= doneChannelPattern =========")
		doneChannelConcurrencyPattern()
	case 8:
		fmt.Println("========= pipelineConcurrencyPattern =========")
		pipelineConcurrencyPattern()
	case 9:
		fmt.Println("========= primesPipeline1_Naive =========")
		primesPipeline1_Naive()
	case 10:
		fmt.Println("========= primesPipeline2FanOut_FanInPattern =========")
		primesPipeline2FanOut_FanInPattern()
	}
}

/*
Go's Fork-Join Model
```mermaid
sequenceDiagram

	participant main
	participant goroutine

	main ->>+ goroutine: fork
	goroutine -->> goroutine: do work
	goroutine -->>- main: join
	note left of main: Join point<br>someone needs to consume

```
*/
func testingGoForkJoinModel() {
	barrierChannel := make(chan string)

	// 1) Executing annonymous functions as goroutine
	go func() {
		barrierChannel <- "Goroutine 1"
	}()

	// Receiving and printing
	fmt.Println(<-barrierChannel)
	close(barrierChannel)

	// ##########################################

	// 2) Wait for the three goroutines to finish
	barrier3Channel := make(chan string, 3)
	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		go func() {
			barrier3Channel <- "Goroutine 2"
			wg.Done()
		}()
		go func() {
			barrier3Channel <- "Goroutine 3"
			wg.Done()
		}()
		go func() {
			barrier3Channel <- "Goroutine 4"
			wg.Done()
		}()

		wg.Wait()
		close(barrier3Channel)
	}()

	// Receiving and printing
	for s := range barrier3Channel {
		fmt.Println(s)
	}

	// ##########################################

	// 3) Wait for the first goroutine to finish
	racer1Channel := make(chan string)
	racer2Channel := make(chan string)
	racer3Channel := make(chan string)

	randMillis := func() time.Duration {
		return time.Duration(rand.Intn(2000-500)+500) * time.Millisecond
	}

	go func() {
		time.Sleep(randMillis())
		racer1Channel <- "🚐"
	}()
	go func() {
		time.Sleep(randMillis())
		racer2Channel <- "🏎️"
	}()
	go func() {
		time.Sleep(randMillis())
		racer3Channel <- "🛺"
	}()

	var winner string
	select {
	case val := <-racer1Channel:
		winner = val
	case val := <-racer2Channel:
		winner = val
	case val := <-racer3Channel:
		winner = val
	}

	fmt.Printf("'%s' WON!!!! 🎉\n", winner)

	// ##########################################

	// 4) At the end
	fmt.Println("Hello World")
}

// Unbuffered channel: send blocks until someone receives.
func blockingChannel() {
	ch := make(chan int) // unbuffered
	ch <- 42             // 🔴 Blocks here because no one is receiving yet (deadlock)
	fmt.Println(<-ch)    // 💤 Never reaches here because the previous line blocks
}

// Buffered channel with size 1: allows sending without blocking while there's space.
// ```go
// ch := make(chan int)    // unbuffered channel, ⚠️ send blocks until a receiver is ready to receive the value
// ch := make(chan int, 1) // buffered channel with capacity 1, ⚠️ allows one value to be sent without blocking (the buffer holds that value)
// ```
func nonBlockingChannel() {
	ch := make(chan int, 1) // buffered channel of size 1
	ch <- 42                // ✅ Does not block because buffer has space
	fmt.Println(<-ch)       // ✅ Prints 42, consumes the value
}

// Unbuffered channel but sending in a goroutine: avoids blocking the main goroutine.
func blockingChannelWithGoRoutine() {
	ch := make(chan int) // unbuffered channel

	go func() {
		fmt.Println("goroutine sending")
		ch <- 42 // 🔴 Blocks here until someone receives
		fmt.Println("goroutine unblocked")
	}()

	fmt.Println("main consuming")
	fmt.Println(<-ch) // ✅ Receives the value, unblocks the goroutine
	fmt.Println("main end")
}

// 1) Send a char slice chars to a buffered charChannel.
// 2) Consume charChannel using a for loop.
// 3) Close the chanel once all chars are on the channel, causing the
// for loop to consume.
/*```go
select {
case charChannel <- s:
	// sent successfully
default:
	// channel is full, skip or handle overflow
}
// ---- VS ----
// blocking sending
charChannel <- s
```
*/
func forSelectPattern() {
	// buffered channel, queue capacity: 3
	// FIFO
	charChannel := make(chan string, 3)
	chars := []string{"a", "b", "c"}

	for _, c := range chars {
		fmt.Printf("charChannel <- '%s'\n", c)
		select {
		case charChannel <- c:
		default:
			fmt.Println("Channel is full, handle overflow")
		}
	}

	fmt.Println("Before closing channel")
	close(charChannel)
	fmt.Println("After closing channel")

	// Blocking loop, will start after the channel is closed
	for c := range charChannel {
		fmt.Printf("'%s' <- charChannel\n", c)
	}

	fmt.Println("After consuming all messages")
}

// asyncCharGoroutine2 demonstrates concurrent producer-consumer behavior using channels.
//
// - A goroutine sends 3 characters ("a", "b", "c") to a buffered channel with a 500ms delay between sends.
// - The main goroutine consumes values from the channel using `for range`, also sleeping 500ms after each receive.
// - Both sender and receiver run in parallel, overlapping their delays.
// - The total execution time is ~1500ms, not ~3000ms, due to this concurrency.
// ```go
//
//	// Blocking for, will start until the first message is sent to myChannelñ
//	for msg := range myChannel {
//	   ...
//	}
//
// // Will block until the channel is closed
//
// ```
// Expected Output (timing may vary slightly):
//
//	sender: sending 'a'
//	receiver: got 'a'
//	sender: sending 'b'
//	receiver: got 'b'
//	sender: sending 'c'
//	receiver: got 'c'
//	receiver: done consuming, duration: ~1500 ms
func asyncCharGoroutine2() {
	start := time.Now()

	// Create a buffered channel with capacity for 3 items.
	charChannel := make(chan string, 3)

	// Start a concurrent sender.
	go func() {
		chars := []string{"a", "b", "c"}
		for _, c := range chars {
			fmt.Printf("sender: sending '%s'\n", c)
			charChannel <- c
			time.Sleep(500 * time.Millisecond)
		}
		fmt.Printf("closing charChannel\n")
		// close(charChannel) // Try commenting and the for will never end
	}()

	fmt.Printf("receiver: starting to consume\n")
	// Start receiving as soon as data sent to the charChannel (blocking).
	for c := range charChannel {
		fmt.Printf("receiver: got '%s'\n", c)
		time.Sleep(500 * time.Millisecond)
	}
	duration := time.Since(start)
	fmt.Printf("receiver: done consuming, duration: %v ms\n", duration.Milliseconds())
}

// Allow parent goroutine to cancel children
func doneChannelConcurrencyPattern() {

	done := make(chan bool)

	go doWork(done)

	time.Sleep(time.Second * 3)

	fmt.Println("closing done channel")
	close(done) // stops doWork loop
}

// done <-chan bool: passing as a read only channel
func doWork(done <-chan bool) {
	for {
		select {
		case <-done:
			fmt.Println("done channel call received")
			return
		default:
			fmt.Println("DOING WORK")
		}
	}
}

// Pipeline, passing (and transforming data) in stages
// This example is synchronous
// Start --input--> stage1(input) --d1--> stage2(d1) --d2--> End
// Execute with: go run main.go 4 | grep pipelineConcurrencyPattern
func pipelineConcurrencyPattern() {
	// Start
	input := []int{2, 3, 4, 7, 1}

	// stage1
	d1Channel := sliceToChannel(input)

	// stage2
	d2Channel := square(d1Channel)

	// stage3
	output := make([]int, 0) // Empty slice, same as "var output []int"
	fmt.Printf("pipelineConcurrencyPattern: starting for consume\n")
	for n := range d2Channel {
		fmt.Printf("pipelineConcurrencyPattern: consume %d\n", n)
		output = append(output, n)
	}
	fmt.Printf("pipelineConcurrencyPattern: collected: %v\n", output)

	// End
}

// Returns read-only int channel
func sliceToChannel(input []int) <-chan int {
	out := make(chan int)
	go func() {
		for _, n := range input {
			fmt.Printf("sliceToChannel goroutine: out <- %d\n", n)
			out <- n
			// Puts "n" in the channel and then blocks itself
		}
		// fmt.Printf("sliceToChannel goroutine: closing\n")
		close(out) // makes the for loop in "square" end
	}()
	// returned even if the nested goroutine hasn't finished
	// fmt.Printf("sliceToChannel: returning out channel\n")
	return out
}

// Receives and returns read only channel
func square(input <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		// We can iterate in "input" even if its capacity is 1
		// Wait until there is a value in "input" channel
		for n := range input {
			fmt.Printf("square goroutine: out <- (%d * %d)\n", n, n)
			out <- n * n
			// Puts "n * n" in the channel and then blocks itself
		}
		// fmt.Printf("square goroutine: closing\n")
		close(out)
	}()
	// returned even if the nested goroutine hasn't finished
	// fmt.Printf("square: returning out channel\n")
	return out
}

// Pipeline
/*
```mermaid
graph TB

	Start((Start))
	End((End))

	Start --> numberGenerator
  	numberGenerator --channel--> FilterPrimes
	FilterPrimes --channel--> TakeN
	TakeN --channel--> End
```
*/
func primesPipeline1_Naive() {
	start := time.Now()

	done := make(chan int)
	defer close(done)

	randomNumberFetcher := func() int { return rand.Intn(500000000) }
	randomIntStream := generatorFunc(done, randomNumberFetcher)
	primeStream := primeFinder(done, randomIntStream)
	for prime := range take(done, primeStream, 10) {
		fmt.Println(prime)
	}

	fmt.Println(time.Since(start))
}

// Pipeline
// FanOut-FanIn pattern
// https://www.youtube.com/watch?v=wELNUHb3kuA
/*
```mermaid
graph TB

	Start((Start))
	End((End))

	Start --> numberGenerator

	subgraph FanOutStage
	    numberGenerator --channel--> FilterPrimes1[FilterPrimes#1]
	    numberGenerator --channel--> FilterPrimes2[FilterPrimes#2]
	    numberGenerator --channel--> FilterPrimes3[FilterPrimes#3]
		numberGenerator --channel--> FilterPrimes4[FilterPrimes#CPUCount]
	end
	subgraph FanInStage
	    FilterPrimes1 --channel--> FanIn
	    FilterPrimes2 --channel--> FanIn
	    FilterPrimes3 --channel--> FanIn
		FilterPrimes4 --channel--> FanIn
	end

	FanIn --channel--> TakeN

	TakeN --channel--> End
```
*/
func primesPipeline2FanOut_FanInPattern() {
	start := time.Now()

	done := make(chan int)
	defer close(done)

	randomNumberGenerator := func() int { return rand.Intn(500000000) }
	randomIntStream := generatorFunc(done, randomNumberGenerator)

	// fan out
	CPUCount := runtime.NumCPU()
	fmt.Printf("=> FanOut-FanIn Pattern running on %d CPUS\n", CPUCount)
	primeFinderChannelSlice := make([]<-chan int, CPUCount)
	for i := range CPUCount {
		primeFinderChannelSlice[i] = primeFinder(done, randomIntStream)
	}

	// fan in
	fannedInStream := fanIn(done, primeFinderChannelSlice...)
	for prime := range take(done, fannedInStream, 10) {
		fmt.Println(prime)
	}
	fmt.Println(time.Since(start))
}

// This function will send one value at a time (unbuffered channel).
// The next function in the chain will unblock [generatorFunc] causing to generate a new value.
// `T any, K any` - Generics
// `done <-chan K` - channel of type K
// `fn func() T` - function to call on each iteration
// `<-chan T` - returns a readonly channel of type T
func generatorFunc[T any, K any](done <-chan K, fn func() T) <-chan T {
	stream := make(chan T) // Unbuffered, synchronous channel
	// Spawn goroutine
	go func() {
		defer close(stream) // Close the stream when the goroutine ends
		for {
			select {
			case <-done:
				return
			case stream <- fn():
				// for-select pattern, while not `done`, i.e. channel is not closed,
				// call `fn` and send the data to the channel
			}
		}
	}()
	return stream
}

func take[T any, K any](done <-chan K, stream <-chan T, n int) <-chan T {
	taken := make(chan T)
	go func() {
		defer close(taken)
		for i := 0; i < n; i++ {
			select {
			case <-done:
				return
			case taken <- <-stream:
				// unblocks `stream` channel (unbuffered), so we can receive the next value
				// `<-stream` - take a value from `stream`
				// `taken <- value` - and pass it to `taken`
			}
		}
	}()
	return taken
}

// Slow pipeline stage
func primeFinder(done <-chan int, randIntStream <-chan int) <-chan int {
	isPrime := func(randomInt int) bool {
		for i := randomInt - 1; i > 1; i-- {
			if randomInt%i == 0 {
				return false
			}
		}
		return true
	}

	primes := make(chan int)

	go func() {
		defer close(primes)
		for {
			select {
			case <-done:
				return
			case randInt := <-randIntStream:
				// fmt.Printf("primeFinder %d\n", randInt)
				if isPrime(randInt) {
					primes <- randInt
				}
			}
		}
	}()

	return primes
}

// Analogy - Coffee shop, a queue of clients wants coffee.
// wg.Add(1)  ➕☕ A new customer enters the coffee shop and wants a coffee.
// wg.Done()  ✅☕ One customer got their coffee and leaves.
// wg.Wait()  ⏳🧑‍🍳 The barista waits until all customers have received their coffee before closing the shop.
func fanIn[T any](done <-chan int, channels ...<-chan T) <-chan T {
	var wg sync.WaitGroup
	fannedInStream := make(chan T)

	/// Transfer data from input channel to fannedInStream channel
	transfer := func(c <-chan T) {
		defer wg.Done()
		for i := range c {
			select {
			case <-done:
				return
			case fannedInStream <- i:
			}
		}
	}

	for _, c := range channels {
		wg.Add(1)
		go transfer(c)
	}

	// Waits for all the transfers to finish
	go func() {
		wg.Wait()
		close(fannedInStream)
	}()
	return fannedInStream
}
