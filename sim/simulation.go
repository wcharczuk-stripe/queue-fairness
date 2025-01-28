package sim

import (
	"fmt"
	"math"
	"math/rand/v2"
	"time"
)

type Simulation struct {
	Config     SimulationConfig
	Clock      Clock
	TaskQueue  TaskQueue
	Workers    WorkerLookup
	RandSource rand.Source

	r *rand.Rand
}

func (s *Simulation) Init() {
	if s.RandSource == nil {
		s.RandSource = rand.NewPCG(rand.Uint64(), rand.Uint64())
	}
	if s.Clock == nil {
		s.Clock = NewSimulatedClock(time.Now())
	}
	if s.TaskQueue == nil {
		s.TaskQueue = NewSimpleTaskQueue()
	}
	s.r = rand.New(s.RandSource)
	s.Workers = s.generateWorkers()
}

func (s *Simulation) Simulate() SimulationResults {
	startTime := s.Clock.Now()
	var lastTimestamp, displayLastTimestamp, currentTimestamp time.Time = startTime, startTime, startTime
	var resultsByBucket resultsByBucket

	var resultState = &results{
		tasks: make([]*Task, 0, 1_000_000),
	}
	for { // hot loop
		currentTimestamp = s.Clock.Now()
		if currentTimestamp.Sub(startTime) > s.Config.DurationOrDefault() {
			break
		}
		s.simulateTick(currentTimestamp, currentTimestamp.Sub(lastTimestamp), resultState)
		s.Clock.Wait(s.Config.TickIntervalOrDefault())
		lastTimestamp = currentTimestamp
		if displayLastTimestamp.IsZero() {
			displayLastTimestamp = currentTimestamp
		} else if currentTimestamp.Sub(displayLastTimestamp) >= s.Config.ResultsBucketingIntervalOrDefault() {
			log(
				fmt.Sprintf("closing results bucket (by interval %v)", s.Config.ResultsBucketingIntervalOrDefault()),
				logTag{"ts", currentTimestamp.Format("15:04")},
				logTag{"elapsed", currentTimestamp.Sub(startTime)},
				logTag{"tql", s.TaskQueue.Len()},
				logTag{"ctp", len(resultState.tasks)},
			)
			resultsByBucket = append(resultsByBucket, resultState)
			displayLastTimestamp = currentTimestamp
			resultState = &results{
				tasks: make([]*Task, 0, 1_000_000),
			}
		}
	}
	return s.processResults(currentTimestamp, resultsByBucket)
}

func (s *Simulation) generateWorkers() WorkerLookup {
	output := make(WorkerLookup)
	for x := 0; x < s.Config.WorkerCountOrDefault(); x++ {
		output.Add(&Worker{
			ID:       x,
			Tasks:    make(TaskLookup, s.Config.WorkerTaskSlotsOrDefault()),
			MaxTasks: s.Config.WorkerTaskSlotsOrDefault(),
		})
	}
	return output
}

func (s *Simulation) simulateTick(currentTimestamp time.Time, elapsedSinceLastTick time.Duration, state *results) {
	s.tickTaskArrivals(currentTimestamp, elapsedSinceLastTick)
	s.tickWorkerPoll(currentTimestamp)
	s.tickWorkerComplete(currentTimestamp, state)
}

func (s *Simulation) tickTaskArrivals(currentTimestamp time.Time, elapsedSinceLastTick time.Duration) {
	newTaskCount := s.randomNewTaskCount(elapsedSinceLastTick)
	for x := 0; x < newTaskCount; x++ {
		t := Task{
			ID:           NewUUID(),
			CreatedUTC:   currentTimestamp,
			Priority:     s.randomPriority(),
			WorkDuration: s.randomWorkDuration(),
		}
		t.FairnessKey, t.Fairness = s.randomFairness()
		s.TaskQueue.Push(t)
	}
}

func (s *Simulation) tickWorkerPoll(currentTimestamp time.Time) {
	for _, w := range s.Workers {
		for len(w.Tasks) < w.MaxTasks {
			t, ok := s.TaskQueue.Pull()
			if !ok {
				return
			}
			t.DispatchedUTC = currentTimestamp
			w.Tasks.Add(t)
		}
	}
}

func (s *Simulation) tickWorkerComplete(currentTimestamp time.Time, state *results) {
	for _, w := range s.Workers {
		var completed []*Task
		for _, t := range w.Tasks {
			if currentTimestamp.Sub(t.DispatchedUTC) >= t.WorkDuration {
				completed = append(completed, t)
			}
		}
		for _, t := range completed {
			t.CompletedUTC = currentTimestamp
			w.Tasks.Del(t)
			state.push(t)
		}
	}
}

func (s *Simulation) randomFairness() (fairnessKey string, fairness float64) {
	if len(s.Config.FairnessWeights) == 0 {
		return "", 1.0
	}
	fairnessKey = s.randomFairnessKey()
	fairness = s.Config.FairnessWeights[fairnessKey]
	return
}

func (s *Simulation) randomFairnessKey() string {
	if len(s.Config.FairnessKeyWeights) == 0 {
		return ""
	}
	return RandomKeyByWeight(s.r, s.Config.FairnessKeyWeights)
}

func (s *Simulation) randomPriority() Priority {
	if len(s.Config.PriorityWeights) == 0 {
		return P2
	}
	return RandomKeyByWeight(s.r, s.Config.PriorityWeights)
}

func (s *Simulation) randomNewTaskCount(elapsedSinceLastTick time.Duration) int {
	elapsedSeconds := float64(elapsedSinceLastTick) / float64(time.Second)
	tasksMean := math.Floor(float64(s.Config.TasksPerSecondOrDefault()) * elapsedSeconds)
	return int(s.randomNormal(tasksMean, tasksMean))
}

func (s *Simulation) randomWorkDuration() time.Duration {
	return time.Duration(RandomNormal(s.r, float64(s.Config.TaskDurationMeanOrDefault()), float64(s.Config.TaskDurationStdDevOrDefault())))
}

func (s *Simulation) randomNormal(desiredMean, desiredStdDev float64) float64 {
	return RandomNormal(s.r, desiredMean, desiredStdDev)
}

type resultsByBucket []*results

type results struct {
	tasks []*Task
}

func (r *results) push(t *Task) {
	r.tasks = append(r.tasks, t)
}
