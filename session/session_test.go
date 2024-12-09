package session

import "testing"

func TestSessionEmptySession(t *testing.T) {
	s := PomoSession{}

	if s.Status() != PomoSessionWork {
		t.Fatalf("Initial state is not work")
	}
}

func TestSessionWorkAfterShortBreak(t *testing.T) {
	s := PomoSession{
		status: PomoSessionShortBreak,
	}

	if s.GetNextStatus() != PomoSessionWork {
		t.Fatalf("Expected work after short break")
	}
}

func TestSessionWorkAfterLongBreak(t *testing.T) {
	s := PomoSession{
		status: PomoSessionLongBreak,
	}

	if s.GetNextStatus() != PomoSessionWork {
		t.Fatalf("Expected work after short break")
	}
}

func genStates(workSessionsBreak int) []PomoSessionStatus {
	N := (workSessionsBreak * 2)
	l := make([]PomoSessionStatus, N)

	l[0] = PomoSessionWork

	for i := 1; i < (N - 1); i += 2 {
		l[i] = PomoSessionShortBreak
		l[i+1] = PomoSessionWork
	}

	l[N-1] = PomoSessionLongBreak
	return l
}

func TestSessionLoop(t *testing.T) {
	N_ITERATIONS := 50
	var N_WORK_SESSIONS int = 4

	s := PomoSession{
		status:            PomoSessionWork,
		WorkSessionsBreak: N_WORK_SESSIONS,
		workedSessions:    0,
	}
	states := genStates(N_WORK_SESSIONS)

	expecetedNWorks := 0
	for i := 0; i < N_ITERATIONS; i++ {

		sessSt := s.Status()
		expectedSt := states[i%len(states)]

		if sessSt != expectedSt {
			t.Fatal(
				"iteration: ", i,
				"status: ", sessSt,
				"expected: ", expectedSt,
			)
		}

		nextSt := s.GetNextStatus()

		if nextSt == PomoSessionWork {
			expecetedNWorks++
		}

		s.SetNextStatus(nextSt)

		nWorks := s.CompletedWorkSessions()
		if expecetedNWorks != nWorks {
			t.Fatal(
				"Iteration: ", i,
				"worked sessions: ", nWorks,
				"expected: ", expecetedNWorks,
			)
		}
	}

}
