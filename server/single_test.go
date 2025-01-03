package server

import (
	"errors"
	"net"
	"net/http"
	"net/rpc"
	"os"
	pomoController "pomogo/controller"
	pomoSession "pomogo/session"
	pomoTimer "pomogo/timer"
	"testing"
	"time"
)

// ========
// FIXTURES
// ========

var zeroDurationCfg = pomoSession.SessionStateDurationConfig{
	PomoSessionWork:       time.Duration(0),
	PomoSessionShortBreak: time.Duration(0),
	PomoSessionLongBreak:  time.Duration(0),
}

var zeroDurationFactory = zeroDurationCfg.GetDurationFactory()

func sessionFactory() pomoSession.PomoSessionIface {
	nWorkSessions := 4
	return &pomoSession.PomoSession{
		WorkSessionsBreak: nWorkSessions,
	}
}

func ssContainerFactory() *pomoController.SingleControllerContainer {
	contFact := func() pomoController.PomoControllerIface {
		ctrl, _ := pomoController.ControllerFactory(
			pomoController.PomoControllerSessionOpt(sessionFactory),
			pomoController.PomoControllerTimerOpt(func() pomoTimer.PomoTimerIface {
				return new(pomoTimer.MockCbTimer)
			}),
			pomoController.PomoControllerDurationF(func() pomoSession.SessionStateDurationFactory {
				return zeroDurationFactory
			}),
		)
		return ctrl
	}

	return &pomoController.SingleControllerContainer{
		ControllerFactory: contFact,
	}
}

// =====
// TESTS
// =====

// Run the test to start-stop
func TestSSStartstop(t *testing.T) {

	SOCKET := "/tmp/pomo_test.sock"

	onListen := func(l net.Listener, s *rpc.Server) error {
		errCh := make(chan error)
		go func() {
			if err := http.Serve(l, s); err != nil {
				errCh <- err
			}
		}()

		// VERY MUCH BRUTE FORCE...
		select {
		case err, ok := <-errCh:
			if !ok {
				t.Fatal("Closed error ch")
			}
			t.Fatal(err)
		case <-time.After(time.Second):
		}

		ssC, err := SingleSessionClientFactory(
			SingleClientRpcHttpConnect("unix", SOCKET),
		)

		if err != nil {
			t.Fatal(err)
		}

		st, err := ssC.Play()

		if err != nil {
			t.Fatal(err)
		}

		if st.State == pomoController.PomoControllerStopped {
			t.Fatal("Expected non stopped status")
		}

		return nil
	}

	serv, err := SingleSessionServerFactory(
		SingleServerContainerOpt(
			ssContainerFactory,
		),
	)

	// OPTION LIKE BUT USED OUTSIDE THE CONTEXT OF A FUNCTION.
	doListen := SingleServerRpcUnixRegOpt(
		SOCKET,
		rpc.NewServer,
		onListen,
	)

	if err != nil {
		t.Fatal(err)
	}

	undoListen, err := doListen(serv)

	if err != nil {
		t.Fatal(err)
	}

	_, err = undoListen(serv)

	if err != nil {
		t.Fatal(err)
	}

	// CHECK THAT THE SOCKET FILE HAS BEEN CLEANED UP.
	if _, err := os.Stat(SOCKET); !errors.Is(err, os.ErrNotExist) {
		t.Fatal(err)
	}
}

// TODO: INCLUDE MORE TESTS. TEST ERRORS AND COMBINATION.
