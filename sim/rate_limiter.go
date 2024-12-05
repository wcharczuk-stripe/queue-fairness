package sim

import "time"

func NewRateLimiter(c Clock, limitActions uint32, limitQuantum time.Duration) RateLimiter {
	return &rateLimiter{
		clock:        c,
		limitActions: limitActions,
		limitQuantum: limitQuantum,
		lastUpdate:   c.Now(),
		tokens:       float64(limitActions) * float64(limitQuantum/time.Second),
	}
}

// RateLimiter is a type with an [RateLimiter.Allow] function
// returns a bool based on how often the function is called.
type RateLimiter interface {
	Allow() bool
	Commit()
}

type rateLimiter struct {
	clock        Clock
	limitActions uint32
	limitQuantum time.Duration
	lastUpdate   time.Time
	tokens       float64 // action seconds
}

func (rl *rateLimiter) Allow() bool {
	now := rl.clock.Now()
	defer func() {
		rl.lastUpdate = now
	}()

	elapsed := now.Sub(rl.lastUpdate)
	elapsedSeconds := float64(elapsed) / float64(time.Second)
	rl.tokens += float64(rl.limitActions) * elapsedSeconds
	if rl.tokens > float64(rl.limitActions) {
		rl.tokens = float64(rl.limitActions)
	}
	return rl.tokens > 1.0
}

func (rl *rateLimiter) Commit() {
	rl.tokens -= 1.0
}
