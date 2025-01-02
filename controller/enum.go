package controller

import (
	"encoding/json"
	"fmt"
	"pomogo/session"
)

// ===================
// PomoControllerState
// ===================

type PomoControllerState int

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

func (s *PomoControllerState) UnmarshalJSON(b []byte) error {

	fmt.Println("PomoControllerState unmarshal")
	sr := string(b)
	switch sr {
	case "Work":
		*s = PomoControllerWork
	case "ShortBreak":
		*s = PomoControllerShortBreak
	case "LongBreak":
		*s = PomoControllerLongBreak
	case "Paused":
		*s = PomoControllerPause
	case "Stopped":
		*s = PomoControllerStopped
	default:
		return &json.MarshalerError{}
	}

	return nil
}

func (s *PomoControllerState) MarshalJSON() ([]byte, error) {

	fmt.Println("PomoControllerState marshal")
	var sr string
	switch *s {
	case PomoControllerWork:
		sr = `"Work"`
	case PomoControllerShortBreak:
		sr = `"ShortBreak"`
	case PomoControllerLongBreak:
		sr = `"LongBreak"`
	case PomoControllerPause:
		sr = `"Paused"`
	case PomoControllerStopped:
		sr = `"Stopped"`
	default:
		return nil, &json.MarshalerError{}
	}
	return []byte(sr), nil
}

const (
	PomoControllerWork PomoControllerState = iota
	PomoControllerShortBreak
	PomoControllerLongBreak
	PomoControllerPause
	PomoControllerStopped
)

// =======================
// PomoControllerEventType
// =======================

type PomoControllerEventType int

const (
	PomoControllerEventTypePlay PomoControllerEventType = iota
	PomoControllerEventTypeStop
	PomoControllerEventTypePause
	PomoControllerEventTypeNextState
)

func (s PomoControllerEventType) String() string {

	switch s {
	case PomoControllerEventTypePlay:
		return "Play"
	case PomoControllerEventTypeStop:
		return "Stop"
	case PomoControllerEventTypePause:
		return "Pause"
	case PomoControllerEventTypeNextState:
		return "NextState"
	}

	panic("Impossible PomoControllerEventType value")
}
