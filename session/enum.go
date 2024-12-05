package session

type PomoSessionStatus int8

const (
	PomoSessionWork PomoSessionStatus = iota
	PomoSessionShortBreak
	PomoSessionLongBreak
)

func (s PomoSessionStatus) String() string {
	switch s {
	case PomoSessionWork:
		return "Work"
	case PomoSessionShortBreak:
		return "ShortBreak:"
	case PomoSessionLongBreak:
		return "LongBreak"
	}

	panic("Impossible session status operator")
}
