package common

import (
	"log"
	"testing"
	"time"
)

func dummyFunc(input chan int, limiter *ConcurrentLimiter) {
	for true {
		data := <-input
		if data == 0 {
			break
		}
		go func(data int) {
			limiter.Begin()
			log.Printf("processing %v\n", data)
			time.Sleep(300 * time.Millisecond)
			log.Printf("processing %v done\n", data)
			limiter.End()
		}(data)
	}
}

func TestConcurrentLimiter(t *testing.T) {
	maxConcurrency := 3
	limiter := NewConcurrentLimiter(maxConcurrency)
	input := make(chan int)

	go dummyFunc(input, limiter)

	for i := 1; i < 20; i++ {
		input <- i
	}
	input <- 0

	for i := 0; i < 3; i++ {
		log.Printf("concurrency: %v\n", limiter.GetCurrentSize())
		time.Sleep(100 * time.Millisecond)
	}
}
