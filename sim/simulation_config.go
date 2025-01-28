package sim

import "time"

// SimulationConfig are parameters to the simulation.
type SimulationConfig struct {
	Duration                 time.Duration
	TickInterval             time.Duration
	ResultsBucketingInterval time.Duration
	TaskDurationMean         time.Duration
	TaskDurationStdDev       time.Duration

	PriorityWeights    map[Priority]int
	FairnessKeyWeights map[string]int
	FairnessWeights    map[string]float64

	WorkerCount     int
	WorkerTaskSlots int
	TasksPerSecond  int
}

func (sc SimulationConfig) DurationOrDefault() time.Duration {
	if sc.Duration > 0 {
		return sc.Duration
	}
	return 1 * time.Hour
}

func (sc SimulationConfig) TickIntervalOrDefault() time.Duration {
	if sc.TickInterval > 0 {
		return sc.TickInterval
	}
	return 500 * time.Millisecond
}

func (sc SimulationConfig) ResultsBucketingIntervalOrDefault() time.Duration {
	if sc.ResultsBucketingInterval > 0 {
		return sc.ResultsBucketingInterval
	}
	return 15 * time.Minute
}

func (sc SimulationConfig) TaskDurationMeanOrDefault() time.Duration {
	if sc.TaskDurationMean > 0 {
		return sc.TaskDurationMean
	}
	return 100 * time.Millisecond
}

func (sc SimulationConfig) TaskDurationStdDevOrDefault() time.Duration {
	if sc.TaskDurationStdDev > 0 {
		return sc.TaskDurationStdDev
	}
	return 50 * time.Millisecond
}

func (sc SimulationConfig) WorkerCountOrDefault() int {
	if sc.WorkerCount > 0 {
		return sc.WorkerCount
	}
	return 32
}

func (sc SimulationConfig) WorkerTaskSlotsOrDefault() int {
	if sc.WorkerTaskSlots > 0 {
		return sc.WorkerTaskSlots
	}
	return 100
}

func (sc SimulationConfig) TasksPerSecondOrDefault() int {
	if sc.TasksPerSecond > 0 {
		return sc.TasksPerSecond
	}
	return 3200
}
