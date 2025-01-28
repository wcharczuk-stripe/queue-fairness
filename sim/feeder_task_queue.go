package sim

import (
	"math/rand/v2"
	"time"
)

type Limit struct {
	Actions uint32
	Quantum time.Duration
}

// NewFeederTaskQueue returns a new feeder task queue with a given set of settings.
//
// The "feeder" task queue type tries to honor absolute rate limits across the fairness keys
// as given by rates formed by counts and a duration of time (e.g. 1000 events per second).
//
// The net effect of this is strictly there is an upperbound to throughput for the system, which
// may be desirable to limit the impact of bursts on downstream systems.
func NewFeederTaskQueue(r *rand.Rand, c Clock, rateLimitsByFairnessKey map[string]Limit) TaskQueue {
	rateLimiters := make(map[string]RateLimiter, len(rateLimitsByFairnessKey))
	for key, lim := range rateLimitsByFairnessKey {
		rateLimiters[key] = NewRateLimiter(c, lim.Actions, lim.Quantum)
	}
	return &feederTaskQueue{
		fairnessKeyRateLimiters: rateLimiters,
		storage:                 make(map[string]*Queue[*Task]),
		r:                       r,
	}
}

type feederTaskQueue struct {
	len                     int
	storage                 map[string]*Queue[*Task]
	fairnessKeyRateLimiters map[string]RateLimiter
	r                       *rand.Rand
}

func (q *feederTaskQueue) Len() int {
	return q.len
}

func (q *feederTaskQueue) Push(t Task) {
	if q.storage[t.FairnessKey] == nil {
		q.storage[t.FairnessKey] = &Queue[*Task]{}
	}
	q.storage[t.FairnessKey].Push(&t)
	q.len++
}

func (q *feederTaskQueue) Pull() (task *Task, ok bool) {
	key, ok := q.getKey()
	if !ok {
		return
	}
	if rl, ok := q.fairnessKeyRateLimiters[key]; ok {
		rl.Commit()
	}
	task, ok = q.storage[key].Pop()
	q.len--
	return
}

func (q *feederTaskQueue) getKey() (key string, ok bool) {
	for fairnessKey := range q.storage {
		if q.storage[fairnessKey].Len() == 0 {
			continue
		}
		if rl, ok := q.fairnessKeyRateLimiters[fairnessKey]; ok && !rl.Allow() {
			continue
		}
		key = fairnessKey
		ok = true
		return
	}
	return
}
