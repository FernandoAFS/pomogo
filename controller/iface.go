package controller

import "time"

type PomoControllerIface interface {
	Status() PomoControllerStatus
	Pause(now time.Time) bool
	Play(now time.Time) bool
	Skip(now time.Time) bool
	Stop(now time.Time) bool
}

// ==========
// STATUS MSG
// ==========

type PomoControllerStatus struct {
	State          PomoControllerState
	TimeLeft       *time.Duration
	PausedAt       *time.Time
	WorkedSessions uint
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
