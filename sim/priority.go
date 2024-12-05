package sim

type Priority int

func (p Priority) String() string {
	switch p {
	case P0:
		return "P0"
	case P1:
		return "P1"
	case P2:
		return "P2"
	case P3:
		return "P3"
	case P4:
		return "P4"
	default:
		return ""
	}
}

const (
	P0 = iota
	P1
	P2
	P3
	P4
)

const (
	DefaultPriority = P2
)
