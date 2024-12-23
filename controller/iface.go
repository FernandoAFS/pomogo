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

// ==========
// STATUS MSG
// ==========

type PomoControllerStatus struct {
	State          PomoControllerState
	TimeLeft       *time.Duration
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
