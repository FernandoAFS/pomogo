// Interfaces and data structures used in for multiple implementations

package controller

import "time"

type PomoControllerIface interface {
	Status() PomoControllerStatus
	Pause(now time.Time) error
	Play(now time.Time) error
	Skip(now time.Time) error
	Stop(now time.Time) error
}

// Manages lifecycle of controller object.
type PomoControllerContainerIface interface {
	GetController() PomoControllerIface
	RemoveController()
}

// ======
// STATUS
// ======

type PomoControllerStatus struct {
	State          PomoControllerState
	TimeLeft       *StatusDuration
	PausedAt       *time.Time
	WorkedSessions int
}

// ======
// EVENTS
// ======

type PomoControllerEventArgsPlay struct {
	At                   time.Time
	CurrentState         PomoControllerState
	NextState            PomoControllerState
	CurrentStateDuration time.Duration
}

type PomoControllerEventArgsStop struct {
	At           time.Time
	CurrentState PomoControllerState
	TimeSpent    time.Duration
	TimeLeft     time.Duration
}

type PomoControllerEventArgsPause struct {
	At           time.Time
	CurrentState PomoControllerState
	TimeSpent    time.Duration
	TimeLeft     time.Duration
}

type PomoControllerEventArgsNextState struct {
	At           time.Time
	CurrentState PomoControllerState
	NextState    PomoControllerState
	TimeLeft     time.Duration
}

// ===========
// STATUS TIME
// ===========

type StatusDuration time.Duration

func (d *StatusDuration) UnmarshalJSON(b []byte) error {

	strB := string(b)

	in_, err := time.ParseDuration(strB[1 : len(strB)-1])
	if err != nil {
		return err
	}
	*d = StatusDuration(in_)
	return nil
}

func (s *StatusDuration) MarshalJSON() ([]byte, error) {
	td := time.Duration(*s).String()
	return []byte(`"` + td + `"`), nil
}
