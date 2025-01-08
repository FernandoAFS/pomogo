package controller

import (
	pomoSession "github.com/FernandoAFS/pomogo/session"
	pomoTimer "github.com/FernandoAFS/pomogo/timer"
	"sync"
	"time"
)

// =================================
// DEFAULT CONTROLLER IMPLEMENTATION
// =================================

// Main business logic component.
// Includes session for sequential state management, timer for background state change.
// It also sends events.
// pausedAt and end-of-state are for pause and status data.
type PomoController struct {
	session         pomoSession.PomoSessionIface
	timer           pomoTimer.PomoTimerIface
	durationFactory pomoSession.SessionStateDurationFactory

	errorSink      func(err error)
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

// Optional error event wrapper
func (c *PomoController) errorEvent(err error) {
	if c.errorSink == nil {
		return
	}
	c.errorSink(err)
}

// Optional play event wrapper
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

// Optional play event wrapper
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

// Freeze timer in time
func (c *PomoController) Pause(now time.Time) error {
	c.locker.Lock()
	defer c.locker.Unlock()

	if c.pauseAt != nil {
		c.errorEvent(ErrPausedTimer)
		return ErrPausedTimer
	}
	c.pauseAt = &now

	if err := c.timer.Cancel(); err != nil {
		c.errorEvent(err)
		return err
	}
	c.pauseEvent(now)
	return nil
}

// Start paused timer or resume paused timer
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

	c.errorEvent(ErrRunningTimer)
	return ErrRunningTimer
}

// Run on a paused timer
func (c *PomoController) resume(now time.Time) error {

	stateTimeLeft := c.endOfState.Sub(*c.pauseAt)
	then := now.Add(stateTimeLeft)

	cb := func() {
		nextStatus := c.session.GetNextStatus()
		if err := c.runTimer(then, nextStatus); err != nil {
			c.errorEvent(err)
		}
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

// call at the end of state timer event
func (c *PomoController) nextTimer(now time.Time) error {
	// RECURSIVE CALLING THE TIMER MUST BE DONE IN A THREAD SAFE WAY.
	c.locker.Lock()
	defer c.locker.Unlock()

	// THIS TWO CHECK SHOULDN'T HAPPEN. IF THERE IS ANY EDGE RACE CONDITION
	// IT'S BETTER TO STOP HERE.

	if c.pauseAt != nil {
		c.errorEvent(ErrPausedTimer)
		return ErrPausedTimer
	}

	if c.endOfState == nil {
		c.errorEvent(ErrStoppedTimer)
		return ErrStoppedTimer
	}

	nextStatus := c.session.GetNextStatus()
	if err := c.runTimer(now, nextStatus); err != nil {
		return err
	}

	// Fire only next status event if the next state is coming.
	c.nextStateEvent(now)
	return nil
}

// start waiting for next timer event.
func (c *PomoController) runTimer(now time.Time, status pomoSession.PomoSessionStatus) error {
	statusDuration := c.durationFactory(status)
	then := now.Add(statusDuration)

	cb := func() {
		if err := c.nextTimer(then); err != nil {
			c.errorEvent(err)
		}
	}

	if err := c.timer.WaitCb(statusDuration, cb); err != nil {
		c.errorEvent(err)
		return err
	}

	c.session.SetNextStatus(status)
	eos := now.Add(statusDuration)
	c.endOfState = &eos
	return nil
}

// Jump to the next status inmediately
func (c *PomoController) Skip(now time.Time) error {
	c.locker.Lock()
	defer c.locker.Unlock()

	// It's ok to skip a paused timer but it will start the next timer right
	// away

	if c.endOfState == nil {
		c.errorEvent(ErrStoppedTimer)
		return ErrStoppedTimer
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

// Reset controller to initial status
func (c *PomoController) Stop(now time.Time) error {
	c.locker.Lock()
	defer c.locker.Unlock()

	// THIS MUST DISMISS EVERY RUNNING GOROUTINE.
	if c.endOfState == nil {
		c.errorEvent(ErrStoppedTimer)
		return ErrStoppedTimer
	}

	if err := c.timer.Cancel(); err != nil {
		c.errorEvent(err)
		return err
	}

	c.stopEvent(now)
	c.endOfState = nil
	return nil
}
