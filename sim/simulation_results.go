package sim

import "time"

type SimulationResults struct {
	TasksProcessed int
	ElapsedTime    time.Duration

	CountByPriority    map[Priority]int
	CountByFairnessKey map[string]int

	QueuedAvg time.Duration
	QueuedP95 time.Duration

	QueuedAvgByPriority map[Priority]time.Duration
	QueuedP95ByPriority map[Priority]time.Duration

	QueuedAvgByFairnessKey map[string]time.Duration
	QueuedP95ByFairnessKey map[string]time.Duration
}

func (s *Simulation) processResults(finalTimestamp time.Time, state resultsByBucket) (res SimulationResults) {
	res.CountByPriority = make(map[Priority]int)
	res.CountByFairnessKey = make(map[string]int)

	res.QueuedAvgByPriority = make(map[Priority]time.Duration)
	res.QueuedP95ByPriority = make(map[Priority]time.Duration)

	res.QueuedAvgByFairnessKey = make(map[string]time.Duration)
	res.QueuedP95ByFairnessKey = make(map[string]time.Duration)

	allQueued := []time.Duration{}
	queuedByPriority := make(map[Priority][]time.Duration)
	queuedByFairnessKey := make(map[string][]time.Duration)

	for _, hour := range state {
		res.TasksProcessed += len(hour.tasks)
		for _, t := range hour.tasks {
			queued := t.DispatchedUTC.Sub(t.CreatedUTC)
			allQueued = append(allQueued, queued)
			res.CountByPriority[t.Priority]++
			res.CountByFairnessKey[t.FairnessKey]++
			queuedByPriority[t.Priority] = append(queuedByPriority[t.Priority], queued)
			queuedByFairnessKey[t.FairnessKey] = append(queuedByFairnessKey[t.FairnessKey], queued)
		}
	}
	for s.TaskQueue.Len() > 0 {
		t, ok := s.TaskQueue.Pull()
		if !ok {
			break
		}
		queued := finalTimestamp.Sub(t.CreatedUTC)
		allQueued = append(allQueued, queued)
		res.CountByPriority[t.Priority]++
		queuedByPriority[t.Priority] = append(queuedByPriority[t.Priority], queued)
		queuedByFairnessKey[t.FairnessKey] = append(queuedByFairnessKey[t.FairnessKey], queued)
	}
	for p, times := range queuedByPriority {
		res.QueuedAvgByPriority[p] = AvgDurations(times)
		res.QueuedP95ByPriority[p] = p95(times)
	}
	for key, times := range queuedByFairnessKey {
		res.QueuedAvgByFairnessKey[key] = AvgDurations(times)
		res.QueuedP95ByFairnessKey[key] = p95(times)
	}
	res.QueuedAvg = AvgDurations(allQueued)
	res.QueuedP95 = p95(allQueued)
	return
}
