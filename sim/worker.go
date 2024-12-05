package sim

type WorkerLookup = Lookup[int, *Worker]

type Worker struct {
	ID       int
	MaxTasks int
	Tasks    TaskLookup
}

func (w *Worker) Key() int { return w.ID }
