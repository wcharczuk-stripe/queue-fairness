package sim

import "time"

// SimulationConfig are parameters to the simulation.
type SimulationConfig struct {
	Duration           time.Duration
	TickInterval       time.Duration
	CompactionInterval time.Duration
	WorkerCount        int
	WorkerTaskSlots    int
	TasksPerSecond     int
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

func (sc SimulationConfig) CompactionIntervalOrDefault() time.Duration {
	if sc.CompactionInterval > 0 {
		return sc.CompactionInterval
	}
	return 15 * time.Minute
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