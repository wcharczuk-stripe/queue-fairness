package sim

// NewPriorityQueue returns a new priority queue.
//
// A priority queue is very similar to a heap but
// instead of bring intrinsically ordered, you can supply a
// separate priority.
func NewPriorityQueue[T any]() *PriorityQueue[T] {
	return &PriorityQueue[T]{
		heap: &Heap[PriorityQueueItem[T]]{
			LessFn: func(i, j PriorityQueueItem[T]) bool {
				return i.Priority < j.Priority
			},
		},
	}
}

// PriorityQueue is a heap of items with priorities.
type PriorityQueue[T any] struct {
	heap *Heap[PriorityQueueItem[T]]
}

// PriorityQueueItem is an item in the priority queue.
type PriorityQueueItem[T any] struct {
	Item     T
	Priority int
}

// Len returns the length of the priority queue.
func (pq *PriorityQueue[T]) Len() int {
	return pq.heap.Len()
}

// Push pushes an item into the priority queue.
func (pq *PriorityQueue[T]) Push(item T, priority int) {
	pq.heap.Push(PriorityQueueItem[T]{
		Item:     item,
		Priority: priority,
	})
}

// Peek returns the minimum item and its priority but does not remove it.
func (pq *PriorityQueue[T]) Peek() (item T, priority int, ok bool) {
	var pi PriorityQueueItem[T]
	pi, ok = pq.heap.Peek()
	if !ok {
		return
	}
	item = pi.Item
	priority = pi.Priority
	return
}

// Pop pops an item off the priority queue.
func (pq *PriorityQueue[T]) Pop() (item T, priority int, ok bool) {
	var pi PriorityQueueItem[T]
	pi, ok = pq.heap.Pop()
	if !ok {
		return
	}
	item = pi.Item
	priority = pi.Priority
	return
}

// Push pushes an item into the priority queue.
func (pq *PriorityQueue[T]) PushPop(inputItem T, inputPriority int) (outputItem T, outputPriority int, ok bool) {
	var pi PriorityQueueItem[T]
	pi, ok = pq.heap.PushPop(PriorityQueueItem[T]{
		Item:     inputItem,
		Priority: inputPriority,
	})
	if !ok {
		return
	}
	outputItem = pi.Item
	outputPriority = pi.Priority
	return
}
