package controller

import (
	"fmt"
	pomoSession "pomogo/session"
	pomoTimer "pomogo/timer"
	"testing"
	"time"
)

// ========
// FIXTURES
// ========

func sessionFactory() *pomoSession.PomoSession {
	var nWorkSessions uint = 4
	return &pomoSession.PomoSession{
		WorkSessionsBreak: nWorkSessions,
	}
}

func mockControllerFactory(
	timer pomoTimer.PomoTimerIface,
	session pomoSession.PomoSessionIface,
	playEventSink func(status PomoControllerEventArgsPlay),
	stopEventSink func(status PomoControllerEventArgsStop),
	pauseEventSink func(status PomoControllerEventArgsPause),
	nextStateEventSink func(event PomoControllerEventArgsNextState),
) *PomoController {

	sessDurationCfg := pomoSession.SessionStateDurationConfig{
		PomoSessionWork:       time.Duration(0),
		PomoSessionShortBreak: time.Duration(0),
		PomoSessionLongBreak:  time.Duration(0),
	}

	durFactory := sessDurationCfg.GetDurationFactory()

	return &PomoController{
		session:            session,
		timer:              timer,
		durationFactory:    durFactory,
		playEventSink:      playEventSink,
		stopEventSink:      stopEventSink,
		pauseEventSink:     pauseEventSink,
		nextStateEventSink: nextStateEventSink,
	}
}

// ============
// ACTION TESTS
// ============

func TestControllerRunStop(t *testing.T) {

	refNow := time.Date(2024, 12, 04, 0, 0, 0, 0, time.UTC)
	timer := pomoTimer.MockTimerFactory(refNow)
	session := sessionFactory()
	controller := mockControllerFactory(timer, session, nil, nil, nil, nil)

	if st := controller.Status().State; st != PomoControllerStopped {
		t.Fatalf("Controller state is %s instead of working", st)
	}

	if err := controller.Play(refNow); err != nil {
		t.Fatal(err)
	}

	if st := controller.Status().State; st != PomoControllerWork {
		t.Fatalf("Controller state is %s instead of working", st)
	}

	if err := controller.Stop(refNow); err != nil {
		t.Fatal(err)
	}

	if st := controller.Status().State; st != PomoControllerStopped {
		t.Fatalf("Controller state is %s instead of stop", st)
	}

}

func TestControllerRunPause(t *testing.T) {

	refNow := time.Date(2024, 12, 04, 0, 0, 0, 0, time.UTC)
	timer := pomoTimer.MockTimerFactory(refNow)
	session := sessionFactory()
	controller := mockControllerFactory(timer, session, nil, nil, nil, nil)

	if st := controller.Status().State; st != PomoControllerStopped {
		t.Fatalf("Controller state is %s instead of working", st)
	}

	if err := controller.Play(refNow); err != nil {
		t.Fatal(err)
	}

	if st := controller.Status().State; st != PomoControllerWork {
		t.Fatalf("Controller state is %s instead of working", st)
	}

	if err := controller.Pause(refNow); err != nil {
		t.Fatal(err)
	}

	if st := controller.Status().State; st != PomoControllerPause {
		t.Fatalf("Controller state is %s instead of pause", st)
	}

	if err := controller.Play(refNow); err != nil {
		t.Fatal(err)
	}

}

func TestControllerNextState(t *testing.T) {

	refNow := time.Date(2024, 12, 04, 0, 0, 0, 0, time.UTC)
	timer := pomoTimer.MockTimerFactory(refNow)
	session := sessionFactory()
	controller := mockControllerFactory(timer, session, nil, nil, nil, nil)

	if st := controller.Status().State; st != PomoControllerStopped {
		t.Fatalf("Controller state is %s instead of working", st)
	}

	if err := controller.Play(refNow); err != nil {
		t.Fatal(err)
	}

	if st := controller.Status().State; st != PomoControllerWork {
		t.Fatalf("Controller state is %s instead of working", st)
	}

	N_ITER := 50
	for i := 0; i < N_ITER; i++ {
		if controller.pauseAt != nil {
			t.Fatalf("Controller is paused at iteration %d", i)
		}

		if controller.endOfState == nil {
			t.Fatalf("Controller is stopped at iteration %d", i)
		}

		if st := controller.Status().State; st == PomoControllerStopped {
			t.Fatalf("Controller is stopped status at iteration %d", i)
		}

		if st := controller.Status().State; st == PomoControllerPause {
			t.Fatalf("Controller is pause status at iteration %d", i)
		}

		sessionSt := SessionToControllerState(session.Status())
		controllerSt := controller.Status().State

		if sessionSt != controllerSt {
			t.Fatalf(
				"Unexpected state controller: %s, session: %s at iteration %d",
				sessionSt,
				controllerSt,
				i,
			)
		}

		timer.ForceDone()
	}
}

func TestControllerNextStatePause(t *testing.T) {

	refNow := time.Date(2024, 12, 04, 0, 0, 0, 0, time.UTC)
	timer := pomoTimer.MockTimerFactory(refNow)
	session := sessionFactory()
	controller := mockControllerFactory(timer, session, nil, nil, nil, nil)

	if st := controller.Status().State; st != PomoControllerStopped {
		t.Fatalf("Controller state is %s instead of working", st)
	}

	if err := controller.Play(refNow); err != nil {
		t.Fatal(err)
	}

	if st := controller.Status().State; st != PomoControllerWork {
		t.Fatalf("Controller state is %s instead of working", st)
	}

	N_ITER := 50
	for i := 0; i < N_ITER; i++ {

		if controller.pauseAt != nil {
			t.Fatalf("Controller is paused at iteration %d", i)
		}

		if controller.endOfState == nil {
			t.Fatalf("Controller is stopped at iteration %d", i)
		}

		if st := controller.Status().State; st == PomoControllerStopped {
			t.Fatalf("Controller is stopped status at iteration %d", i)
		}

		if st := controller.Status().State; st == PomoControllerPause {
			t.Fatalf("Controller is pause status at iteration %d", i)
		}

		if err := controller.Pause(refNow); err != nil {
			t.Fatal(err)
		}

		if st := controller.Status().State; st != PomoControllerPause {
			t.Fatalf("Controller is not pause status at iteration %d", i)
		}

		if err := controller.Play(refNow); err != nil {
			t.Fatal(err)
		}

		sessionSt := SessionToControllerState(session.Status())
		controllerSt := controller.Status().State

		if sessionSt != controllerSt {
			t.Fatalf(
				"Unexpected state controller: %s, session: %s at iteration %d",
				sessionSt,
				controllerSt,
				i,
			)
		}

		timer.ForceDone()
	}
}

// ===========
// EVENT TESTS
// ===========

func TestControllerPlayEvent(t *testing.T) {

	eventTime := time.Date(2024, 12, 06, 0, 0, 0, 0, time.UTC)
	refNow := time.Date(2024, 12, 04, 0, 0, 0, 0, time.UTC)

	eventPlayed := false
	playSink := func(event PomoControllerEventArgsPlay) {
		if event.At != eventTime {
			t.Fatalf(
				"Event time (%s), is not expected (%s)",
				event.At,
				eventTime,
			)
		}

		if event.CurrentState != PomoControllerWork {
			t.Fatalf(
				"Expected event state is Work and got %s",
				event.CurrentState,
			)
		}
		eventPlayed = true
	}

	timer := pomoTimer.MockTimerFactory(refNow)
	session := sessionFactory()
	controller := mockControllerFactory(timer, session, playSink, nil, nil, nil)

	if st := controller.Status().State; st != PomoControllerStopped {
		t.Fatalf("Controller state is %s instead of working", st)
	}

	if err := controller.Play(eventTime); err != nil {
		t.Fatal(err)
	}

	if !eventPlayed {
		t.Fatalf("Event not runned after play action")
	}

	if st := controller.Status().State; st != PomoControllerWork {
		t.Fatalf(
			"Controller state is %s instead of working",
			st,
		)
	}
}

func TestControllerStopEvent(t *testing.T) {

	eventTime := time.Date(2024, 12, 06, 0, 0, 0, 0, time.UTC)
	refNow := time.Date(2024, 12, 04, 0, 0, 0, 0, time.UTC)

	stopEventDone := false
	stopEventSink := func(event PomoControllerEventArgsStop) {
		if event.At != eventTime {
			t.Fatalf(
				"Event time (%s), is not expected (%s)",
				event.At,
				eventTime,
			)
		}

		if event.CurrentState != PomoControllerWork {
			t.Fatalf(
				"Expected event state is Work and got %s",
				event.CurrentState,
			)
		}

		// TODO: CHECK TIME LEFT AND TIME SPENT...

		stopEventDone = true
	}

	timer := pomoTimer.MockTimerFactory(refNow)
	session := sessionFactory()
	controller := mockControllerFactory(timer, session, nil, stopEventSink, nil, nil)

	if st := controller.Status().State; st != PomoControllerStopped {
		t.Fatalf("Controller state is %s instead of working", st)
	}

	if err := controller.Play(refNow); err != nil {
		t.Fatal(err)
	}

	if st := controller.Status().State; st != PomoControllerWork {
		t.Fatalf(
			"Controller state is %s instead of working",
			st,
		)
	}

	if err := controller.Stop(eventTime); err != nil {
		t.Fatal(err)
	}

	if !stopEventDone {
		t.Fatalf("Stop event must have runned")
	}

	if st := controller.Status().State; st != PomoControllerStopped {
		t.Fatalf("Controller state is %s instead of stop", st)
	}

}

// TODO: INCLUDE TEST FOR: PAUSE/RESUME AND NEXT EVENTS

func TestControllerPauseResumeEvent(t *testing.T) {

	eventTime := time.Date(2024, 12, 06, 0, 0, 0, 0, time.UTC)
	refNow := time.Date(2024, 12, 04, 0, 0, 0, 0, time.UTC)

	playEventPlayed := false
	pauseEventDone := false

	playSink := func(event PomoControllerEventArgsPlay) {
		if event.At != eventTime {
			t.Fatalf(
				"Event time (%s), is not expected (%s)",
				event.At,
				eventTime,
			)
		}

		if event.CurrentState != PomoControllerWork {
			t.Fatalf(
				"Expected event state is Work and got %s",
				event.CurrentState,
			)
		}
		playEventPlayed = true
	}

	pauseEventSink := func(event PomoControllerEventArgsPause) {
		if event.At != eventTime {
			t.Fatalf(
				"Event time (%s), is not expected (%s)",
				event.At,
				eventTime,
			)
		}

		if event.CurrentState != PomoControllerWork {
			t.Fatalf(
				"Expected event state is Work and got %s",
				event.CurrentState,
			)
		}
		pauseEventDone = true
	}

	timer := pomoTimer.MockTimerFactory(refNow)
	session := sessionFactory()
	controller := mockControllerFactory(timer, session, playSink, nil, pauseEventSink, nil)

	if st := controller.Status().State; st != PomoControllerStopped {
		t.Fatalf("Controller state is %s instead of working", st)
	}

	if err := controller.Play(eventTime); err != nil {
		t.Fatal(err)
	}

	if !playEventPlayed {
		t.Fatalf("Play event must have runned")
	}

	if st := controller.Status().State; st != PomoControllerWork {
		t.Fatalf(
			"Controller state is %s instead of working",
			st,
		)
	}

	if err := controller.Pause(eventTime); err != nil {
		t.Fatal(err)
	}

	if !pauseEventDone {
		t.Fatalf("Pause event must have runned")
	}

	if st := controller.Status().State; st != PomoControllerPause {
		t.Fatalf("Controller state is %s instead of pause", st)
	}
}

func TestControllerNextStateEvent(t *testing.T) {

	eventTime := time.Date(2024, 12, 06, 0, 0, 0, 0, time.UTC)
	refNow := time.Date(2024, 12, 04, 0, 0, 0, 0, time.UTC)

	playEventPlayed := false

	playSink := func(event PomoControllerEventArgsPlay) {
		if event.At != eventTime {
			t.Fatalf(
				"Event time (%s), is not expected (%s)",
				event.At,
				eventTime,
			)
		}

		if event.CurrentState != PomoControllerWork {
			t.Fatalf(
				"Expected event state is Work and got %s",
				event.CurrentState,
			)
		}
		playEventPlayed = true
	}

	timer := pomoTimer.MockTimerFactory(refNow)
	session := sessionFactory()

	nextStateEventSinkCounter := 0
	nextStateEventSinkTriggerCounter := 0

	nextStateEvCh := make(chan *string)

	nextStateEventSink := func(event PomoControllerEventArgsNextState) {
		// INCLUDE STATE CHECKING... DOUBLE CHECK WHAT IS RIGHT...

		if event.At != eventTime {
			err := fmt.Sprintf(
				"Next state event time (%s), is not expected (%s)",
				event.At,
				eventTime,
			)
			nextStateEvCh <- &err
		}

		if st := PomoControllerState(session.Status()); event.CurrentState != st {
			err := fmt.Sprintf(
				"Next state event state (%s), is not expected (%s)",
				event.CurrentState,
				st,
			)
			nextStateEvCh <- &err
		}

		if st := PomoControllerState(session.GetNextStatus()); event.NextState != st {
			err := fmt.Sprintf(
				"Next state event state (%s), is not expected (%s)",
				event.CurrentState,
				st,
			)
			nextStateEvCh <- &err
		}

		nextStateEventSinkCounter++
		nextStateEvCh <- nil
	}

	controller := mockControllerFactory(timer, session, playSink, nil, nil, nextStateEventSink)

	// PLAY SETUP
	if st := controller.Status().State; st != PomoControllerStopped {
		t.Fatalf("Controller state is %s instead of working", st)
	}

	if err := controller.Play(eventTime); err != nil {
		t.Fatal(err)
	}

	if !playEventPlayed {
		t.Fatalf("Play event must have runned")
	}

	N_ITERATIONS := 50

	for i := 0; i < N_ITERATIONS; i++ {
		timer.ForceDone()
		nextStateEventSinkTriggerCounter++

		if err := <-nextStateEvCh; err != nil {
			t.Fatalf("Error on iteration %d: %s", i, *err)
		}

		if nextStateEventSinkCounter != nextStateEventSinkTriggerCounter {
			t.Fatalf(
				"Next state event count is %d while expected %d ",
				nextStateEventSinkTriggerCounter,
				nextStateEventSinkCounter,
			)
		}
	}
	close(nextStateEvCh)
}
