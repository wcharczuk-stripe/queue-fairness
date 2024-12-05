package sim

type TaskQueue interface {
	Push(Task)
	Pull() (*Task, bool)
	Len() int
}
