package sim

import "time"

// Clock is a source of time
type Clock interface {
	Now() time.Time
	Wait(time.Duration)
}

type WallClock struct{}

func (wc WallClock) Now() time.Time       { return time.Now() }
func (wc WallClock) Wait(d time.Duration) { time.Sleep(d) }

func NewSimulatedClock(startAt time.Time) Clock {
	return &simulatedClock{ts: startAt}
}

type simulatedClock struct {
	ts time.Time
}

func (sc *simulatedClock) Now() time.Time        { return sc.ts }
func (sc *simulatedClock) Wait(by time.Duration) { sc.ts = sc.ts.Add(by) }
