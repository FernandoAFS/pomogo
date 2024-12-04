package controller

import "time"

type PomoControllerStatus struct{
    State PomoControllerState
    TimeLeft *time.Duration
    PausedAt *time.Time
    WorkedSessions uint
}

type PomoControllerIface interface{
    Status() PomoControllerStatus 
    Pause(now time.Time) bool
    Play(now time.Time) bool
    Skip(now time.Time) bool
    Stop(now time.Time) bool
}
