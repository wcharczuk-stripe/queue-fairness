package sim

import "testing"

func Test_SimpleTaskQueue(t *testing.T) {
	rq := NewSimpleTaskQueue()

	rq.Push(Task{ID: NewUUID()})
	rq.Push(Task{ID: NewUUID()})
	rq.Push(Task{ID: NewUUID()})
	rq.Push(Task{ID: NewUUID()})
	rq.Push(Task{ID: NewUUID()})

	if rq.Len() != 5 {
		t.Errorf("expect tq length to be 5, was %d", rq.Len())
		t.Fail()
	}

	task, ok := rq.Pull()
	if !ok {
		t.Errorf("expect pull to be ok")
		t.Fail()
	}
	if task == nil {
		t.Errorf("expect pull task to not be nil")
		t.Fail()
	}

	task, ok = rq.Pull()
	if !ok {
		t.Errorf("expect pull to be ok")
		t.Fail()
	}
	if task == nil {
		t.Errorf("expect pull task to not be nil")
		t.Fail()
	}

	task, ok = rq.Pull()
	if !ok {
		t.Errorf("expect pull to be ok")
		t.Fail()
	}
	if task == nil {
		t.Errorf("expect pull task to not be nil")
		t.Fail()
	}

	task, ok = rq.Pull()
	if !ok {
		t.Errorf("expect pull to be ok")
		t.Fail()
	}
	if task == nil {
		t.Errorf("expect pull task to not be nil")
		t.Fail()
	}

	task, ok = rq.Pull()
	if !ok {
		t.Errorf("expect pull to be ok")
		t.Fail()
	}
	if task == nil {
		t.Errorf("expect pull task to not be nil")
		t.Fail()
	}

	if rq.Len() != 0 {
		t.Errorf("expect tq length to be 0, was %d", rq.Len())
		t.Fail()
	}
}
