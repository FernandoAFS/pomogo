package session

import "time"

type SessionStateDurationConfig struct {
	PomoSessionWork       time.Duration
	PomoSessionShortBreak time.Duration
	PomoSessionLongBreak  time.Duration
}

func (cfg *SessionStateDurationConfig) GetDurationFactory() SessionStateDurationFactory {

	return func(s PomoSessionStatus) time.Duration {
		switch s {
		case PomoSessionWork:
			return cfg.PomoSessionWork
		case PomoSessionShortBreak:
			return cfg.PomoSessionShortBreak
		case PomoSessionLongBreak:
			return cfg.PomoSessionLongBreak
		}

		panic("Impossible session status operator")
	}

}
