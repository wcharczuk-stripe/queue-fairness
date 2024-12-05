package sim

func NewPrioritySortedTaskQueue() TaskQueue {
	return &prioritySortedTaskQueue{
		storage: NewPriorityQueue[*Task](),
	}
}

type prioritySortedTaskQueue struct {
	storage *PriorityQueue[*Task]
}

func (q *prioritySortedTaskQueue) Len() int {
	return q.storage.Len()
}

func (q *prioritySortedTaskQueue) Push(t Task) {
	q.storage.Push(&t, int(t.Priority))
}

func (q *prioritySortedTaskQueue) Pull() (task *Task, ok bool) {
	task, _, ok = q.storage.Pop()
	return
}
