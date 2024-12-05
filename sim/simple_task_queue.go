package sim

func NewSimpleTaskQueue() TaskQueue {
	return &simpleTaskQueue{
		storage: &Queue[*Task]{},
	}
}

type simpleTaskQueue struct {
	storage *Queue[*Task]
}

func (q *simpleTaskQueue) Len() int {
	return q.storage.Len()
}

func (q *simpleTaskQueue) Push(t Task) {
	q.storage.Push(&t)
}

func (q *simpleTaskQueue) Pull() (task *Task, ok bool) {
	task, ok = q.storage.Pop()
	return
}
