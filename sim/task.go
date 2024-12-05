package sim

import "time"

type TaskLookup = Lookup[UUID, *Task]

type Task struct {
	ID            UUID
	Priority      Priority
	FairnessKey   string
	Fairness      float64
	CreatedUTC    time.Time
	DispatchedUTC time.Time
	CompletedUTC  time.Time
	WorkDuration  time.Duration
}

func (t Task) Key() UUID {
	return t.ID
}
