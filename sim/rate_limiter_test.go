package sim

import (
	"testing"
	"time"
)

func Test_RateLimiter(t *testing.T) {
	c := NewSimulatedClock(time.Now())
	rl := NewRateLimiter(c, 10, time.Second) // 10 actions per second

	for x := 0; x < 9; x++ {
		if !rl.Allow() {
			t.Errorf("Expect the first 9 calls to be true")
			t.FailNow()
		}
		rl.Commit()

		c.Wait(5 * time.Millisecond)
	}

	// should be @ 10 calls in 500 milliseconds
	// and the next call should be debounced
	if rl.Allow() {
		t.Errorf("Expect the 10th call to be false")
		t.FailNow()
	}
	rl.Commit()

	// advance such that we get our rate back
	c.Wait(time.Second)

	for x := 0; x < 9; x++ {
		if !rl.Allow() {
			t.Errorf("Expect the second 9 calls to be true")
			t.FailNow()
		}
		rl.Commit()
		c.Wait(5 * time.Millisecond)
	}

	for x := 0; x < 20; x++ {
		if rl.Allow() {
			t.Errorf("Expect subsequent calls to be false")
			t.FailNow()
		}
		rl.Commit()
		c.Wait(50 * time.Millisecond)
	}
}
