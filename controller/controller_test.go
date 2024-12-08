package controller

import (
	pomoSession "pomogo/session"
	pomoTimer "pomogo/timer"
	"testing"
	"time"
)

var zeroDurationCfg = pomoSession.SessionStateDurationConfig{
	PomoSessionWork:       time.Duration(0),
	PomoSessionShortBreak: time.Duration(0),
	PomoSessionLongBreak:  time.Duration(0),
}

var zeroDurationFactory = zeroDurationCfg.GetDurationFactory()

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
	options ...PomoControllerOption,
) *PomoController {
	cf := PomoControllerFactory{
		session:         session,
		timer:           timer,
		durationFactory: zeroDurationFactory,
	}
	return cf.Create(options...)
}

// ============
// ACTION TESTS
// ============

func TestControllerRunStop(t *testing.T) {

	refNow := time.Date(2024, 12, 04, 0, 0, 0, 0, time.UTC)
	timer := &pomoTimer.MockCbTimer{}
	session := sessionFactory()

	cf := PomoControllerFactory{
		session:         session,
		timer:           timer,
		durationFactory: zeroDurationCfg.GetDurationFactory(),
	}
	controller := cf.Create()

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
	timer := &pomoTimer.MockCbTimer{}
	session := sessionFactory()
	controller := mockControllerFactory(timer, session)

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
	timer := &pomoTimer.MockCbTimer{}
	session := sessionFactory()
	controller := mockControllerFactory(timer, session)

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
	timer := &pomoTimer.MockCbTimer{}
	session := sessionFactory()
	controller := mockControllerFactory(timer, session)

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

	timer := &pomoTimer.MockCbTimer{}
	session := sessionFactory()

	cf := PomoControllerFactory{
		session:         session,
		timer:           timer,
		durationFactory: zeroDurationCfg.GetDurationFactory(),
	}
	controller := cf.Create(PomoControllerOptionPlaySink(playSink))

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

	timer := &pomoTimer.MockCbTimer{}
	session := sessionFactory()

	controller := mockControllerFactory(
		timer,
		session,
		PomoControllerOptionStopSink(stopEventSink),
	)

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

	timer := &pomoTimer.MockCbTimer{}
	session := sessionFactory()

	controller := mockControllerFactory(
		timer,
		session,
		PomoControllerOptionPlaySink(playSink),
		PomoControllerOptionPauseSink(pauseEventSink),
	)

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

	timer := &pomoTimer.MockCbTimer{}
	session := sessionFactory()

	nextStateEventSinkCounter := 0
	nextStateEventSinkTriggerCounter := 0

	nextStateEventSink := func(event PomoControllerEventArgsNextState) {
		// INCLUDE STATE CHECKING... DOUBLE CHECK WHAT IS RIGHT...

		if event.At != eventTime {
			t.Fatalf(
				"Next state event time (%s), is not expected (%s)",
				event.At,
				eventTime,
			)
		}

		if st := PomoControllerState(session.Status()); event.CurrentState != st {
			t.Fatalf(
				"Next state event state (%s), is not expected (%s)",
				event.CurrentState,
				st,
			)
		}

		if st := PomoControllerState(session.GetNextStatus()); event.NextState != st {
			t.Fatalf(
				"Next state event state (%s), is not expected (%s)",
				event.CurrentState,
				st,
			)
		}

		nextStateEventSinkCounter++
	}

	controller := mockControllerFactory(
		timer,
		session,
		PomoControllerOptionPlaySink(playSink),
		PomoControllerOptionNextStateSink(nextStateEventSink),
	)

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

		if nextStateEventSinkCounter != nextStateEventSinkTriggerCounter {
			t.Fatalf(
				"Next state event count is %d while expected %d ",
				nextStateEventSinkTriggerCounter,
				nextStateEventSinkCounter,
			)
		}
	}
}
