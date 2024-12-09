package session

import "time"

type PomoSessionIface interface {
	Status() PomoSessionStatus
	GetNextStatus() PomoSessionStatus
	SetNextStatus(status PomoSessionStatus)
	Reset()
	CompletedWorkSessions() int
}

type SessionStateDurationFactory func(s PomoSessionStatus) time.Duration
