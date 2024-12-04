package controller

import (
	pomoSession "pomogo/session"
	pomoTimer "pomogo/timer"
	"sync"
	"time"
)

type PomoController struct {
	session         pomoSession.PomoSessionIface
	timer           pomoTimer.PomoTimerIface
	durationFactory pomoSession.SessionStateDurationFactory

	pauseAt    *time.Time
	endOfState *time.Time

	locker sync.Mutex
}

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
	timeLeft := c.endOfState.Sub(now)
	return PomoControllerStatus{
		State:          SessionToControllerState(c.session.Status()),
		TimeLeft:       &timeLeft,
		PausedAt:       nil,
		WorkedSessions: workedSessions,
	}
}

func (c *PomoController) Pause(now time.Time) error {
	c.locker.Lock()
	defer c.locker.Unlock()

	if c.pauseAt != nil {
		return PausedTimer
	}
	c.pauseAt = &now

	if err := c.timer.Cancel(); err != nil {
		return err
	}
	return nil
}

func (c *PomoController) Play(now time.Time) error {
	c.locker.Lock()
	defer c.locker.Unlock()

	if c.endOfState == nil {
		c.session.Reset()
		return c.runTimer(now)
	}

	if c.pauseAt != nil {
		return c.resume(now)
	}

	return RunningTimer
}

func (c *PomoController) resume(now time.Time) error {

	stateTimeLeft := c.endOfState.Sub(*c.pauseAt)

	cb := func(then time.Time) {
		nextStatus := c.session.GetNextStatus()
		c.session.SetNextStatus(nextStatus)
		c.runTimer(then)
	}

	if err := c.timer.WaitCb(stateTimeLeft, cb); err != nil {
		return err
	}

	c.pauseAt = nil
	eos := now.Add(stateTimeLeft)
	c.endOfState = &eos
	return nil
}

func (c *PomoController) nextTimer(now time.Time) error {
	// RECURSIVE CALLING THE TIMER MUST BE DONE IN A THREAD SAFE WAY.
	c.locker.Lock()
	defer c.locker.Unlock()

	// THIS TWO CHECK SHOULDN'T HAPPEN. IF THERE IS ANY EDGE RACE CONDITION
	// IT'S BETTER TO STOP HERE.

	if c.pauseAt != nil {
		return PausedTimer
	}

	if c.endOfState == nil {
		return StoppedTimer
	}

	nextStatus := c.session.GetNextStatus()
	c.session.SetNextStatus(nextStatus)
	return c.runTimer(now)
}

// TODO: RENAME THIS FUNCTION...
func (c *PomoController) runTimer(now time.Time) error {
	status := c.session.Status()
	statusDuration := c.durationFactory(status)

	cb := func(then time.Time) { c.nextTimer(then) }
	if err := c.timer.WaitCb(statusDuration, cb); err != nil {
		return err
	}

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
		return StoppedTimer
	}

	if err := c.timer.Cancel(); err != nil {
		return err
	}

	nextStatus := c.session.GetNextStatus()
	c.session.SetNextStatus(nextStatus)
	return c.runTimer(now)
}

func (c *PomoController) Stop(now time.Time) error {
	c.locker.Lock()
	defer c.locker.Unlock()

	// THIS MUST DISMISS EVERY RUNNING GOROUTINE.
	if c.endOfState == nil {
		return StoppedTimer
	}

	if err := c.timer.Cancel(); err != nil {
		return err
	}

	c.endOfState = nil
	return nil
}
