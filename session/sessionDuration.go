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

func DurationFactory(
	PomoSessionWorkDuration time.Duration,
	PomoSessionShortBreakDuration time.Duration,
	PomoSessionLongBreakDuration time.Duration,
) SessionStateDurationFactory {
	return func(s PomoSessionStatus) time.Duration {
		switch s {
		case PomoSessionWork:
			return PomoSessionWorkDuration
		case PomoSessionShortBreak:
			return PomoSessionShortBreakDuration
		case PomoSessionLongBreak:
			return PomoSessionLongBreakDuration
		}
		panic("Impossible session status operator")
	}
}
