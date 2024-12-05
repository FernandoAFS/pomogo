package session

import "testing"

func TestEmptySession(t *testing.T) {
	s := PomoSession{}

	if s.Status() != PomoSessionWork {
		t.Fatalf("Initial state is not work")
	}
}

func TestWorkAfterShortBreak(t *testing.T) {
	s := PomoSession{
		status: PomoSessionShortBreak,
	}

	if s.GetNextStatus() != PomoSessionWork {
		t.Fatalf("Expected work after short break")
	}
}

func TestWorkAfterLongBreak(t *testing.T) {
	s := PomoSession{
		status: PomoSessionLongBreak,
	}

	if s.GetNextStatus() != PomoSessionWork {
		t.Fatalf("Expected work after short break")
	}
}
