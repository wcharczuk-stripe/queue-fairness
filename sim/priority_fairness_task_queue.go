package sim

import "math/rand/v2"

func NewPriorityFairnessTaskQueue(r *rand.Rand) TaskQueue {
	return &priorityFairnessTaskQueue{
		fairnessKeyWeights: make(map[string]float64),
		r:                  r,
	}
}

type priorityFairnessTaskQueue struct {
	len                int
	storage            [5]map[string]map[UUID]*Task
	fairnessKeyWeights map[string]float64
	r                  *rand.Rand
}

func (q *priorityFairnessTaskQueue) Len() int {
	return q.len
}

func (q *priorityFairnessTaskQueue) Push(t Task) {
	if q.storage[t.Priority] == nil {
		q.storage[t.Priority] = make(map[string]map[UUID]*Task)
	}
	if q.storage[t.Priority][t.FairnessKey] == nil {
		q.storage[t.Priority][t.FairnessKey] = make(map[UUID]*Task)
	}
	q.storage[t.Priority][t.FairnessKey][t.ID] = &t
	q.fairnessKeyWeights[t.FairnessKey] = t.Fairness
	q.len++
}

func (q *priorityFairnessTaskQueue) Pull() (task *Task, ok bool) {
	for _, p := range []Priority{P0, P1, P2, P3, P4} {
		if len(q.storage[p]) == 0 {
			continue
		}
		fairnessKey := RandomKeyByWeight(q.r, filterMapBySharedKeys(q.fairnessKeyWeights, q.storage[p]))
		task, ok = mapFirst(q.storage[p][fairnessKey])
		if !ok {
			return
		}
		delete(q.storage[p][fairnessKey], task.ID)
		if len(q.storage[p][fairnessKey]) == 0 {
			delete(q.storage[p], fairnessKey)
		}
		q.len--
		return
	}
	return
}

func filterMapBySharedKeys[K comparable, V0, V1 any](filterThis map[K]V0, byThat map[K]V1) map[K]V0 {
	output := make(map[K]V0)
	for key, value := range filterThis {
		if _, ok := byThat[key]; ok {
			output[key] = value
		}
	}
	return output
}

func mapFirst[K comparable, V any](m map[K]V) (v V, ok bool) {
	for _, v = range m {
		ok = true
		return
	}
	return
}
