package common

import (
	"log"
	"testing"
	"time"
)

func waitToNanosecond(timestamp int64) {
	for {
		if time.Now().UnixNano() >= int64(timestamp) {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func TestRetryLimiter(t *testing.T) {
	InitRetryLimiter(2)

	crtTime := time.Now().UnixNano()
	var limiter *RetryLimiter
	// 1st req in 1st window
	limiter = NewRetryLimiter("Test", 2, 0.5)
	if !limiter.CanRetry() {
		t.Error("retry logic error 0-01")

	}
	if !limiter.CanRetry() {
		t.Error("retry logic error 0-02")
	}
	if limiter.CanRetry() {
		t.Error("retry logic error 0-03")
	}
	// 2nd req in 1st window
	limiter = NewRetryLimiter("Test", 2, 0.5)
	if !limiter.CanRetry() {
		t.Error("retry logic error 0-13")
	}
	if !limiter.CanRetry() {
		t.Error("retry logic error 0-13")
	}
	if limiter.CanRetry() {
		t.Error("retry logic error 0-13")
	}
	log.Println("2nd window starts")
	// 2nd window starts
	waitToNanosecond(crtTime + 500*1000000)

	// 1st req in 2nd window
	limiter = NewRetryLimiter("Test", 2, 0.5)
	if !limiter.CanRetry() {
		t.Error("retry logic error 1-03")
	}
	if !limiter.CanRetry() {
		t.Error("retry logic error 1-03")
	}
	if limiter.CanRetry() {
		t.Error("retry logic error 1-03")
	}
	log.Println(limiter.Debug())

	// 2nd req in 2nd window
	limiter = NewRetryLimiter("Test", 2, 0.5)
	if !limiter.CanRetry() || !limiter.CanRetry() {
		t.Error("retry logic error 1-11")
	}
	if limiter.CanRetry() {
		t.Error("retry logic error 1-12")
	}
	log.Println(limiter.Debug())
}
