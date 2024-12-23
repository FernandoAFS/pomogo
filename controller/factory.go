// FACTORY FOLLOWING FUNCTION OPTIONS PATTERN.

package controller

import (
	pomoSession "pomogo/session"
	pomoTimer "pomogo/timer"
)

// CONSIDER INCLUDING ERRORS IN OPTIONS.

type PomoControllerOption func(*PomoController) PomoControllerOption

// =======
// FACTORY
// =======

type PomoControllerFactory struct {
	Session         pomoSession.PomoSessionIface
	Timer           pomoTimer.PomoTimerIface
	DurationFactory pomoSession.SessionStateDurationFactory
}

func (f *PomoControllerFactory) Create(
	options ...PomoControllerOption,
) *PomoController {

	c := &PomoController{
		session:         f.Session,
		timer:           f.Timer,
		durationFactory: f.DurationFactory,
	}

	for _, opt := range options {
		opt(c)
	}

	return c
}

// =======
// OPTIONS
// =======

func PomoControllerOptionErrorSink(errorSink func(err error)) PomoControllerOption {
	return func(c *PomoController) PomoControllerOption {
		prev := c.errorSink
		c.errorSink = errorSink
		return PomoControllerOptionErrorSink(prev)
	}
}

func PomoControllerOptionPlaySink(
	playEventSink func(event PomoControllerEventArgsPlay),
) PomoControllerOption {
	return func(c *PomoController) PomoControllerOption {
		prev := c.playEventSink
		c.playEventSink = playEventSink
		return PomoControllerOptionPlaySink(prev)
	}
}

func PomoControllerOptionStopSink(
	stopEventSink func(event PomoControllerEventArgsStop),
) PomoControllerOption {
	return func(c *PomoController) PomoControllerOption {
		prev := c.stopEventSink
		c.stopEventSink = stopEventSink
		return PomoControllerOptionStopSink(prev)
	}
}

func PomoControllerOptionPauseSink(
	pauseEventSink func(event PomoControllerEventArgsPause),
) PomoControllerOption {
	return func(c *PomoController) PomoControllerOption {
		prev := c.pauseEventSink
		c.pauseEventSink = pauseEventSink
		return PomoControllerOptionPauseSink(prev)
	}
}

func PomoControllerOptionNextStateSink(
	nextStateEventSink func(event PomoControllerEventArgsNextState),
) PomoControllerOption {
	return func(c *PomoController) PomoControllerOption {
		prev := c.nextStateEventSink
		c.nextStateEventSink = nextStateEventSink
		return PomoControllerOptionNextStateSink(prev)
	}
}
