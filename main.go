package main

import (
	"flag"
	"fmt"
	"math/rand/v2"
	"os"
	"runtime/pprof"
	"sort"
	"time"

	"queue_fairness/sim"
)

var (
	flagRealTime   = flag.Bool("real-time", false, "if we should simulate using real (wall clock) time")
	flagCPUProfile = flag.Bool("cpu-profile", false, "if we should take a cpu profile")
	flagQueueType  = flag.String("queue-type", "feeder", "which queue type to use (simple|priority|fairness|feeder)")

	flagDuration           = flag.Duration("duration", sim.SimulationConfig{}.DurationOrDefault(), "the simulation duration")
	flagCompactionInterval = flag.Duration("compaction-interval", sim.SimulationConfig{}.CompactionIntervalOrDefault(), "the simulation compaction interval")
	flagTickInterval       = flag.Duration("tick-interval", sim.SimulationConfig{}.TickIntervalOrDefault(), "the simulation tick interval")

	flagTaskMean   = flag.Duration("task-mean", sim.SimulationConfig{}.TaskDurationMeanOrDefault(), "the task duration mean")
	flagTaskStdDev = flag.Duration("task-std-dev", sim.SimulationConfig{}.TaskDurationStdDevOrDefault(), "the task duration std dev")
)

func main() {
	flag.Parse()
	s := new(sim.Simulation)
	s.Config.Duration = *flagDuration
	s.Config.CompactionInterval = *flagCompactionInterval
	s.Config.TickInterval = *flagTickInterval
	s.Config.TaskDurationMean = *flagTaskMean
	s.Config.TaskDurationStdDev = *flagTaskStdDev
	s.RandSource = rand.NewPCG(rand.Uint64(), rand.Uint64())

	if *flagRealTime {
		s.Clock = new(sim.WallClock)
	} else {
		s.Clock = sim.NewSimulatedClock(time.Now())
	}

	switch *flagQueueType {
	case "simple":
		s.TaskQueue = sim.NewSimpleTaskQueue()
	case "priority":
		s.TaskQueue = sim.NewPrioritySortedTaskQueue()
	case "fairness":
		s.TaskQueue = sim.NewPriorityFairnessTaskQueue(rand.New(s.RandSource))
	case "feeder":
		s.TaskQueue = sim.NewFeederTaskQueue(rand.New(s.RandSource), s.Clock, map[string]sim.Limit{
			"high":   {Actions: 7000, Quantum: time.Second}, // these mirror 70/20/10 for the fk weights
			"medium": {Actions: 2000, Quantum: time.Second},
			"low":    {Actions: 1000, Quantum: time.Second},
		})
	default:
		fmt.Fprintf(os.Stderr, "invalid queue type: %v\n", *flagQueueType)
		os.Exit(1)
	}

	fmt.Printf("using task queue type:\t\t%v\n", *flagQueueType)
	fmt.Printf("using simulation duration:\t%v\n", s.Config.DurationOrDefault())
	fmt.Printf("using compaction interval:\t%v\n", s.Config.CompactionIntervalOrDefault())
	fmt.Printf("using tick interval:\t\t%v\n", s.Config.TickIntervalOrDefault())
	fmt.Printf("using tasks-per-second:\t\t%v\n", s.Config.TasksPerSecondOrDefault())
	fmt.Printf("using tasks duration mean:\t\t%v\n", s.Config.TaskDurationMeanOrDefault())
	fmt.Printf("using tasks duration std dev:\t\t%v\n", s.Config.TaskDurationStdDevOrDefault())
	fmt.Println()

	s.Init()

	var profileDone func()
	var err error
	if *flagCPUProfile {
		profileDone, err = cpuProfile()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error starting cpu profile: %v\n", err)
			os.Exit(1)
		}
	}

	res := s.Simulate()
	if *flagCPUProfile {
		profileDone()
	}

	fmt.Println()
	fmt.Printf("tasks processed: %d\n", res.TasksProcessed)
	fmt.Printf("queued for \tp95: %v\tavg: %v\n", res.QueuedP95.Round(time.Millisecond).String(), res.QueuedAvg.Round(time.Millisecond).String())
	fmt.Println()
	for _, p := range []sim.Priority{sim.P0, sim.P1, sim.P2, sim.P3, sim.P4} {
		fmt.Printf("queued for by priority %q [%d]\t\tp95: %v\tavg: %v\n",
			p,
			res.CountByPriority[p],
			res.QueuedP95ByPriority[p].Round(time.Millisecond).String(),
			res.QueuedAvgByPriority[p].Round(time.Millisecond).String(),
		)
	}
	for _, key := range sortedKeys(res.QueuedP95ByFairnessKey) {
		fmt.Printf("queued for by fairness key %q [%d]\tp95: %v\tavg: %v\n",
			key,
			res.CountByFairnessKey[key],
			res.QueuedP95ByFairnessKey[key].Round(time.Millisecond).String(),
			res.QueuedAvgByFairnessKey[key].Round(time.Millisecond).String(),
		)
	}
}

func sortedKeys[T any](m map[string]T) (output []string) {
	for key := range m {
		output = append(output, key)
	}
	sort.Strings(output)
	return
}

func cpuProfile() (func(), error) {
	f, err := os.Create("cpu.profile.pb.gz")
	if err != nil {
		return nil, err
	}
	err = pprof.StartCPUProfile(f)
	if err != nil {
		return nil, err
	}
	return func() {
		pprof.StopCPUProfile()
		_ = f.Close()
	}, nil
}
