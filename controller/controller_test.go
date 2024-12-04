package controller

import (
	pomoSession "pomogo/session"
	pomoTimer "pomogo/timer"
	"testing"
	"time"
)

// =============
// DRY FUNCTIONS
// =============

func sessionFactory() *pomoSession.PomoSession {
	var nWorkSessions uint = 4
	return &pomoSession.PomoSession{
		WorkSessionsBreak: nWorkSessions,
	}
}

func mockTestFactory(
	timer pomoTimer.PomoTimerIface,
	session pomoSession.PomoSessionIface,
) *PomoController {

	sessDurationCfg := pomoSession.SessionStateDurationConfig{
		PomoSessionWork:       time.Duration(0),
		PomoSessionShortBreak: time.Duration(0),
		PomoSessionLongBreak:  time.Duration(0),
	}

	durFactory := sessDurationCfg.GetDurationFactory()

	return &PomoController{
		session:         session,
		timer:           timer,
		durationFactory: durFactory,
	}
}

// ==========
// TEST CASES
// ==========

func TestControllerRunStop(t *testing.T) {

	refNow := time.Date(2024, 12, 04, 0, 0, 0, 0, time.UTC)
	timer := pomoTimer.MockTimerFactory(refNow)
	session := sessionFactory()
	controller := mockTestFactory(timer, session)

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
	controller := mockTestFactory(timer, session)

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
	controller := mockTestFactory(timer, session)

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
	controller := mockTestFactory(timer, session)

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

		if err := controller.Pause(refNow); err != nil{
			t.Fatal(err)
		}
		
		if st := controller.Status().State; st != PomoControllerPause {
			t.Fatalf("Controller is not pause status at iteration %d", i)
		}

		if err := controller.Play(refNow); err != nil{
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
