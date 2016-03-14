package common

/*
limiter with same names share throughput counter.

Usage:
limiter = common.NewRetryLimiter("Func1", 2, 0.1)
for limiter.CanRetry() {
    // do something here
}

*/

import (
	"fmt"
	"sync"
	"time"
)

type limiterData struct {
	window         int64
	throughput     int // throughtput in current window
	lastThroughput int // throughtput in last time window
	retryCount     int
	lock           *sync.Mutex
}

var (
	throughputWindow int
	dataMap          map[string]*limiterData
	countChannel     chan string
	dataMapLock      *sync.Mutex
)

func init() {
	throughputWindow = 10
	dataMap = make(map[string]*limiterData)
	countChannel = make(chan string)
	dataMapLock = &sync.Mutex{}

	go func(countChannel chan string) {
		// couunt throughput for all retries in current process
		for {
			writerName, ok := <-countChannel
			if !ok {
				break
			}

			dataMapLock.Lock()
			if _, ok := dataMap[writerName]; !ok {
				crtData := limiterData{
					0, 0, 0, 0, &sync.Mutex{},
				}
				dataMap[writerName] = &crtData
			}
			crtWindow := time.Now().Unix() / int64(throughputWindow)
			crtData := dataMap[writerName]

			if crtWindow > crtData.window {
				// entering new window
				crtData.window = crtWindow
				crtData.lastThroughput = crtData.throughput
				crtData.throughput = 1
			} else {
				crtData.throughput++
			}
			dataMapLock.Unlock()
		}
	}(countChannel)
}

// InitRetryLimiter setup retry limit
func InitRetryLimiter(window int) {
	throughputWindow = window
}

// StopRetryLimiter stop retry counter
func StopRetryLimiter() {
	close(countChannel)
}

// RetryLimiter handler for retry limiter
type RetryLimiter struct {
	name       string
	retryLimit int
	retryCount int
	retryRatio float64
}

// NewRetryLimiter create a new retry limiter, you should create a new retry limiter everytime you calls rpc
func NewRetryLimiter(name string, retryLimit int, retryRatio float64) *RetryLimiter {
	countChannel <- name
	return &RetryLimiter{name, retryLimit, 0, retryRatio}
}

// CanRetry call this function everytime before actuall function call
func (p *RetryLimiter) CanRetry() bool {
	if p.retryCount == 0 {
		// first time always success
		p.retryCount = 1
		return true
	}
	if p.retryCount >= p.retryLimit {
		return false
	}

	dataMapLock.Lock()
	crtData, ok := dataMap[p.name]
	dataMapLock.Unlock()
	if !ok {
		p.retryCount++
		return true
	}
	crtData.lock.Lock()
	defer func() {
		crtData.lock.Unlock()
	}()

	if crtData.lastThroughput >= 0 && float64(crtData.retryCount) > float64(crtData.lastThroughput)*p.retryRatio {
		return false
	}
	p.retryCount++
	return true
}

// Debug for debug only
func (p *RetryLimiter) Debug() string {
	dataMapLock.Lock()
	crtData, ok := dataMap[p.name]
	if !ok {
		dataMapLock.Unlock()
		return fmt.Sprintf("error:%v", dataMap)
	}
	result := fmt.Sprintf("crtTime:%v, win: %v, name: %v, limit:%v, retry: %v, limitRatio: %v, throughput: %v, lastThroughput: %v",
		time.Now().UnixNano(), crtData.window,
		p.name, p.retryLimit, p.retryCount, p.retryRatio, crtData.throughput, crtData.lastThroughput,
	)
	dataMapLock.Unlock()
	return result
}
