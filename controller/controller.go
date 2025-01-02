package controller

import (
	pomoSession "pomogo/session"
	pomoTimer "pomogo/timer"
	"sync"
	"time"
)

// =================================
// DEFAULT CONTROLLER IMPLEMENTATION
// =================================

type PomoController struct {
	session         pomoSession.PomoSessionIface
	timer           pomoTimer.PomoTimerIface
	durationFactory pomoSession.SessionStateDurationFactory

	// SINKS ARE RUN SYNCHRONOUSLY. RUN ASYNC FUNCTIONS INSIDE THE CALLBACK IF
	// NECESSARY.

	// SINKS ARE CONSIDERED OPTIONAL. IF THEY ARE NOT INFORMED THEY WON'T RUN.

	// THIS IS USEFUL SINCE THERE MAY BE BACKGROUND ERRORS
	errorSink func(err error)

	// RUN ON PLAY OR RESUME
	playEventSink  func(event PomoControllerEventArgsPlay)
	stopEventSink  func(event PomoControllerEventArgsStop)
	pauseEventSink func(event PomoControllerEventArgsPause)

	// RUN ON END OF STATE TIME OR ON SKIP STATES
	nextStateEventSink func(event PomoControllerEventArgsNextState)

	pauseAt    *time.Time
	endOfState *time.Time

	locker sync.Mutex
}

// -------------
// GETTER METHOD
// -------------

// Return a status report of the controller
func (c *PomoController) Status() PomoControllerStatus {
	c.locker.Lock()
	defer c.locker.Unlock()

	if c.endOfState == nil {
		return PomoControllerStatus{
			State:          PomoControllerStopped,
			TimeLeft:       nil,
			PausedAt:       nil,
			WorkedSessions: 0,
		}
	}

	workedSessions := c.session.CompletedWorkSessions()

	if c.pauseAt != nil {
		return PomoControllerStatus{
			State:          PomoControllerPause,
			TimeLeft:       nil,
			PausedAt:       c.pauseAt,
			WorkedSessions: workedSessions,
		}
	}

	now := time.Now()
	timeLeft := StatusDuration(c.endOfState.Sub(now))
	return PomoControllerStatus{
		State:          SessionToControllerState(c.session.Status()),
		TimeLeft:       &timeLeft,
		PausedAt:       nil,
		WorkedSessions: workedSessions,
	}
}

// --------------
// EVENT EMITTING
// --------------

func (c *PomoController) errorEvent(err error) {
	if c.errorSink == nil {
		return
	}
	c.errorSink(err)
}

func (c *PomoController) playEvent(now time.Time) {
	if c.playEventSink == nil {
		return
	}

	status := c.session.Status()
	nextStatus := c.session.GetNextStatus()

	playEvent := PomoControllerEventArgsPlay{
		At:                   now,
		CurrentState:         SessionToControllerState(status),
		NextState:            SessionToControllerState(nextStatus),
		CurrentStateDuration: c.durationFactory(status),
	}

	c.playEventSink(playEvent)
}

func (c *PomoController) stopEvent(now time.Time) {
	if c.stopEventSink == nil {
		return
	}

	status := c.session.Status()
	duration := c.durationFactory(status)
	timeLeft := c.endOfState.Sub(now)
	timeSpent := duration - timeLeft

	stopEvent := PomoControllerEventArgsStop{
		At:           now,
		CurrentState: SessionToControllerState(status),
		TimeSpent:    timeSpent,
		TimeLeft:     timeLeft,
	}

	c.stopEventSink(stopEvent)
}

func (c *PomoController) pauseEvent(now time.Time) {
	if c.pauseEventSink == nil {
		return
	}

	status := c.session.Status()
	duration := c.durationFactory(status)
	timeLeft := c.endOfState.Sub(now)
	timeSpent := duration - timeLeft

	pauseEvent := PomoControllerEventArgsPause{
		At:           now,
		CurrentState: SessionToControllerState(status),
		TimeSpent:    timeSpent,
		TimeLeft:     timeLeft,
	}

	c.pauseEventSink(pauseEvent)
}

func (c *PomoController) nextStateEvent(now time.Time) {
	if c.nextStateEventSink == nil {
		return
	}

	status := c.session.Status()
	nextStatus := c.session.GetNextStatus()
	timeLeft := c.endOfState.Sub(now)

	nextStateEvent := PomoControllerEventArgsNextState{
		At:           now,
		CurrentState: SessionToControllerState(status),
		NextState:    SessionToControllerState(nextStatus),
		TimeLeft:     timeLeft,
	}

	c.nextStateEventSink(nextStateEvent)
}

// ------------------
// CONTROLLER ACTIONS
// ------------------

func (c *PomoController) Pause(now time.Time) error {
	c.locker.Lock()
	defer c.locker.Unlock()

	if c.pauseAt != nil {
		c.errorEvent(PausedTimer)
		return PausedTimer
	}
	c.pauseAt = &now

	if err := c.timer.Cancel(); err != nil {
		c.errorEvent(err)
		return err
	}
	c.pauseEvent(now)
	return nil
}

func (c *PomoController) Play(now time.Time) error {
	c.locker.Lock()
	defer c.locker.Unlock()

	if c.endOfState == nil {
		c.session.Reset()
		status := c.session.Status()
		if err := c.runTimer(now, status); err != nil {
			c.errorEvent(err)
			return err
		}
		c.playEvent(now)
		return nil
	}

	if c.pauseAt != nil {
		return c.resume(now)
	}

	c.errorEvent(RunningTimer)
	return RunningTimer
}

func (c *PomoController) resume(now time.Time) error {

	stateTimeLeft := c.endOfState.Sub(*c.pauseAt)
	then := now.Add(stateTimeLeft)

	cb := func() {
		nextStatus := c.session.GetNextStatus()
		c.runTimer(then, nextStatus)
	}

	if err := c.timer.WaitCb(stateTimeLeft, cb); err != nil {
		c.errorEvent(err)
		return err
	}

	c.pauseAt = nil
	eos := now.Add(stateTimeLeft)
	c.endOfState = &eos
	c.playEvent(now)
	return nil
}

func (c *PomoController) nextTimer(now time.Time) error {
	// RECURSIVE CALLING THE TIMER MUST BE DONE IN A THREAD SAFE WAY.
	c.locker.Lock()
	defer c.locker.Unlock()

	// THIS TWO CHECK SHOULDN'T HAPPEN. IF THERE IS ANY EDGE RACE CONDITION
	// IT'S BETTER TO STOP HERE.

	if c.pauseAt != nil {
		c.errorEvent(PausedTimer)
		return PausedTimer
	}

	if c.endOfState == nil {
		c.errorEvent(StoppedTimer)
		return StoppedTimer
	}

	c.nextStateEvent(now)
	nextStatus := c.session.GetNextStatus()
	return c.runTimer(now, nextStatus)
}

// TODO: RENAME THIS FUNCTION...
func (c *PomoController) runTimer(now time.Time, status pomoSession.PomoSessionStatus) error {
	statusDuration := c.durationFactory(status)
	then := now.Add(statusDuration)

	cb := func() { c.nextTimer(then) }
	if err := c.timer.WaitCb(statusDuration, cb); err != nil {
		c.errorEvent(err)
		return err
	}

	c.session.SetNextStatus(status)
	eos := now.Add(statusDuration)
	c.endOfState = &eos
	return nil
}

func (c *PomoController) Skip(now time.Time) error {
	c.locker.Lock()
	defer c.locker.Unlock()

	// IT'S OK TO SKIP A PAUSED TIMER BUT IT WILL START THE NEXT TIMER RIGHT
	// AWAY

	if c.endOfState == nil {
		c.errorEvent(StoppedTimer)
		return StoppedTimer
	}

	if err := c.timer.Cancel(); err != nil {
		c.errorEvent(err)
		return err
	}

	nextStatus := c.session.GetNextStatus()
	c.nextStateEvent(now)
	// This is broken. if error rises it changes the state and keeps the
	// existing work order...
	return c.runTimer(now, nextStatus)
}

func (c *PomoController) Stop(now time.Time) error {
	c.locker.Lock()
	defer c.locker.Unlock()

	// THIS MUST DISMISS EVERY RUNNING GOROUTINE.
	if c.endOfState == nil {
		c.errorEvent(StoppedTimer)
		return StoppedTimer
	}

	if err := c.timer.Cancel(); err != nil {
		c.errorEvent(err)
		return err
	}

	c.stopEvent(now)
	c.endOfState = nil
	return nil
}
