package sim

import (
	"math/rand/v2"
	"testing"
	"time"
)

func Test_FeederTaskQueue(t *testing.T) {
	c := NewSimulatedClock(time.Now())
	r := rand.NewPCG(123, 123)
	rq := NewFeederTaskQueue(rand.New(r), c, map[string]Limit{
		"high":   {Actions: 1000, Quantum: time.Second},
		"medium": {Actions: 500, Quantum: time.Second},
		"low":    {Actions: 100, Quantum: time.Second},
	})

	rq.Push(Task{ID: NewUUID(), FairnessKey: "high", Fairness: 70})
	rq.Push(Task{ID: NewUUID(), FairnessKey: "high", Fairness: 70})
	rq.Push(Task{ID: NewUUID(), FairnessKey: "medium", Fairness: 20})
	rq.Push(Task{ID: NewUUID(), FairnessKey: "medium", Fairness: 20})
	rq.Push(Task{ID: NewUUID(), FairnessKey: "medium", Fairness: 20})
	rq.Push(Task{ID: NewUUID(), FairnessKey: "low", Fairness: 10})
	rq.Push(Task{ID: NewUUID(), FairnessKey: "low", Fairness: 10})
	rq.Push(Task{ID: NewUUID(), FairnessKey: "low", Fairness: 10})
	rq.Push(Task{ID: NewUUID(), FairnessKey: "low", Fairness: 10})

	if rq.Len() != 9 {
		t.Errorf("expect tq length to be 9, was %d", rq.Len())
		t.Fail()
	}
	for x := 0; x < 9; x++ {
		task, ok := rq.Pull()
		if !ok {
			t.Errorf("expect pull to be ok")
			t.Fail()
		}
		if task == nil {
			t.Errorf("expect pull task to not be nil")
			t.Fail()
		}
	}
	if rq.Len() != 0 {
		t.Errorf("expect tq length to be 0, was %d", rq.Len())
		t.Fail()
	}
}
