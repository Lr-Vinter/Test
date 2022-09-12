package limiter

import (
	"sync"
	"time"
)

type Limiter struct {
	mu      sync.Mutex
	counter int

	maxRequest int
	timeline   time.Duration

	timer time.Time
}

func NewLimiter(maxRequest int, timeline time.Duration) *Limiter {
	return &Limiter{
		mu:         sync.Mutex{},
		counter:    0,
		maxRequest: maxRequest,
		timeline:   timeline,
	}
}

func (l *Limiter) CheckSkip() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	timenow := time.Now()
	if timenow.Sub(l.timer) > l.timeline {
		l.timer = time.Now()
		l.counter = 0
	}

	l.counter++
	if l.counter > l.maxRequest {
		return true
	} else {
		return false
	}
}
