package common

// ConcurrentLimiter can limit the concurrent of centain code block
type ConcurrentLimiter struct {
	Enabled bool
	sem     chan bool
}

// NewConcurrentLimiter creates a new limiter
func NewConcurrentLimiter(maxConcurrent int) *ConcurrentLimiter {
	sem := make(chan bool, maxConcurrent)
	if maxConcurrent <= 0 {
		return &ConcurrentLimiter{
			false, sem,
		}
	}
	return &ConcurrentLimiter{
		true, sem,
	}
}

// Begin call this function before the actuall logic
func (p *ConcurrentLimiter) Begin() bool {
	if p.Enabled {
		select {
		case p.sem <- true:
			return true
		default:
			return false
		}
	}
	return true
}

// End release a concurrent lock
func (p *ConcurrentLimiter) End() bool {
	if p.Enabled {
		<-p.sem
	}
	return true
}

// GetCurrentSize returns how many go routines is current running
func (p *ConcurrentLimiter) GetCurrentSize() int {
	return len(p.sem)
}

// GetMaxCocurrent gets the capacity of the limiter
func (p *ConcurrentLimiter) GetMaxCocurrent() int {
	return cap(p.sem)
}

// IsFull returns if cocurrentcy limit is reached
func (p *ConcurrentLimiter) IsFull() bool {
	if !p.Enabled {
		return false
	}
	return len(p.sem) >= cap(p.sem)
}
