// FACTORY FOLLOWING FUNCTION OPTIONS PATTERN.

package controller

import (
	pomoSession "pomogo/session"
	pomoTimer "pomogo/timer"
)

type PomoControllerOption func(*PomoController) PomoControllerOption

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

type PomoControllerFactory struct {
	session         pomoSession.PomoSessionIface
	timer           pomoTimer.PomoTimerIface
	durationFactory pomoSession.SessionStateDurationFactory
}

func (f *PomoControllerFactory) Create(
	options ...PomoControllerOption,
) *PomoController {

	c := &PomoController{
		session:         f.session,
		timer:           f.timer,
		durationFactory: f.durationFactory,
	}

	for _, opt := range options {
		opt(c)
	}

	return c
}
