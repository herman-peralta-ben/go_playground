// go run main.go

package main

import (
	"fmt"
	"time"
)

func main() {
	// This will be called when main finishes, even with error
	defer func() {
		fmt.Println("finished main")
	}()

	for val := range takeStream(tickerStream(200), 5) {
		fmt.Println(val)
	}
	chainedStreamOperators()
	chainedStreamOperators2()
}

func chainedStreamOperators() {
	src := tickerStream(200)

	mapped := mapStream(src, func(i int) int {
		return i * 2
	})

	filtered := filterStream(mapped, func(i int) bool {
		return i%4 == 0
	})

	for val := range filtered {
		fmt.Println("chainedStreamOperators:", val)
		if val > 20 {
			break
		}
	}
}

func chainedStreamOperators2() {
	for val := range filterStream(
		mapStream(
			tickerStream(200), func(i int) int {
				return i * 2
			}), func(i int) bool {
			return i%4 == 0
		}) {
		fmt.Println("chainedStreamOperators2:", val)
		if val > 20 {
			break
		}
	}
}

func tickerStream(millis int) <-chan int {
	ch := make(chan int)
	go func() {
		defer close(ch)
		ticker := time.NewTicker(time.Duration(millis) * time.Millisecond)
		defer ticker.Stop()

		i := 0
		for {
			select {
			case <-ticker.C:
				ch <- i
				i++
			}
		}
	}()
	return ch
}

func mapStream(in <-chan int, f func(int) int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for val := range in {
			out <- f(val)
		}
	}()
	return out
}

func filterStream(in <-chan int, pred func(int) bool) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for val := range in {
			if pred(val) {
				out <- val
			}
		}
	}()
	return out
}

func takeStream(in <-chan int, n int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		i := 0
		for val := range in {
			if i >= n {
				break
			}
			out <- val
			i++
		}
	}()
	return out
}
