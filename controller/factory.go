// FACTORY FOLLOWING FUNCTION OPTIONS PATTERN.

package controller

import (
	pomoSession "github.com/FernandoAFS/pomogo/session"
	pomoTimer "github.com/FernandoAFS/pomogo/timer"
)

// CONSIDER INCLUDING ERRORS IN OPTIONS.

type PomoControllerOption func(*PomoController) (PomoControllerOption, error)

// =======
// FACTORY
// =======

// Aggregator of Pomodoro Controller options and initializes pointer.
func ControllerFactory(
	options ...PomoControllerOption,
) (*PomoController, error) {
	c := new(PomoController)
	for _, opt := range options {
		_, err := opt(c)
		if err != nil {
			return nil, err
		}

	}
	return c, nil
}

// =======
// OPTIONS
// =======

// Sets controller session from factory
func PomoControllerSessionOpt(sessionF func() pomoSession.PomoSessionIface) PomoControllerOption {
	return func(c *PomoController) (PomoControllerOption, error) {
		prev := c.session
		c.session = sessionF()
		return func(c *PomoController) (PomoControllerOption, error) {
			c.session = prev
			return PomoControllerSessionOpt(sessionF), nil
		}, nil
	}
}

// Sets controller timer from factory
func PomoControllerTimerOpt(timerF func() pomoTimer.PomoTimerIface) PomoControllerOption {
	return func(c *PomoController) (PomoControllerOption, error) {
		prev := c.timer
		c.timer = timerF()
		return func(c *PomoController) (PomoControllerOption, error) {
			c.timer = prev
			return PomoControllerTimerOpt(timerF), nil
		}, nil
	}
}

// Sets controller duration factory from factory
func PomoControllerDurationF(
	durationF func() pomoSession.SessionStateDurationFactory,
) PomoControllerOption {
	return func(c *PomoController) (PomoControllerOption, error) {
		prev := c.durationFactory
		c.durationFactory = durationF()
		return func(c *PomoController) (PomoControllerOption, error) {
			c.durationFactory = prev
			return PomoControllerDurationF(durationF), nil
		}, nil
	}
}

// Sets error sinks
func PomoControllerOptionErrorSink(errorSink func(err error)) PomoControllerOption {
	return func(c *PomoController) (PomoControllerOption, error) {
		prev := c.errorSink
		c.errorSink = errorSink
		return PomoControllerOptionErrorSink(prev), nil
	}
}

// Sets play sinks
func PomoControllerOptionPlaySink(
	playEventSink func(event PomoControllerEventArgsPlay),
) PomoControllerOption {
	return func(c *PomoController) (PomoControllerOption, error) {
		prev := c.playEventSink
		c.playEventSink = playEventSink
		return PomoControllerOptionPlaySink(prev), nil
	}
}

// Sets stop sinks
func PomoControllerOptionStopSink(
	stopEventSink func(event PomoControllerEventArgsStop),
) PomoControllerOption {
	return func(c *PomoController) (PomoControllerOption, error) {
		prev := c.stopEventSink
		c.stopEventSink = stopEventSink
		return PomoControllerOptionStopSink(prev), nil
	}
}

// Sets pause sinks
func PomoControllerOptionPauseSink(
	pauseEventSink func(event PomoControllerEventArgsPause),
) PomoControllerOption {
	return func(c *PomoController) (PomoControllerOption, error) {
		prev := c.pauseEventSink
		c.pauseEventSink = pauseEventSink
		return PomoControllerOptionPauseSink(prev), nil
	}
}

// Sets next state sinks
func PomoControllerOptionNextStateSink(
	endOfStateEventSink func(event PomoControllerEventArgsNextState),
) PomoControllerOption {
	return func(c *PomoController) (PomoControllerOption, error) {
		prev := c.endOfStateEventSink
		c.endOfStateEventSink = endOfStateEventSink
		return PomoControllerOptionNextStateSink(prev), nil
	}
}

// Create an event listener that runs command on every event
func PomoControllerHook(command string) PomoControllerOption {
	// Check that the command is reacheble and executable

	return func(c *PomoController) (PomoControllerOption, error) {

		prevPlay := c.playEventSink
		prevStop := c.stopEventSink
		prevPause := c.pauseEventSink
		prevNe := c.endOfStateEventSink
		prevErr := c.errorSink

		c.playEventSink = PlayExecHook(command)
		c.stopEventSink = StopExecHook(command)
		c.pauseEventSink = PauseExecHook(command)
		c.endOfStateEventSink = NextStateExecHook(command)
		c.errorSink = ErrorExecHook(command)

		return func(c *PomoController) (PomoControllerOption, error) {
			c.playEventSink = prevPlay
			c.stopEventSink = prevStop
			c.pauseEventSink = prevPause
			c.endOfStateEventSink = prevNe
			c.errorSink = prevErr

			return PomoControllerHook(command), nil
		}, nil
	}
}
