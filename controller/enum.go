package controller

import "pomogo/session"

type PomoControllerState uint8

func (s PomoControllerState) String() string {

	switch s {
	case PomoControllerWork:
		return "Work"
	case PomoControllerShortBreak:
		return "ShortBreak"
	case PomoControllerLongBreak:
		return "LongBreak"
	case PomoControllerPause:
		return "Paused"
	case PomoControllerStopped:
		return "Stopped"
	}

	panic("Impossible PomoControllerState value")
}

func SessionToControllerState(s session.PomoSessionStatus) PomoControllerState {

	switch s {
	case session.PomoSessionWork:
		return PomoControllerWork
	case session.PomoSessionShortBreak:
		return PomoControllerShortBreak
	case session.PomoSessionLongBreak:
		return PomoControllerLongBreak
	}

	panic("Impossible PomoSessionStatus value")
}

const (
	PomoControllerWork PomoControllerState = iota
	PomoControllerShortBreak
	PomoControllerLongBreak
	PomoControllerPause
	PomoControllerStopped
)
