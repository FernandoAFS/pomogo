package timer

import "time"

// TIMER WITH LOCK PROTECTION
type MockCbTimer struct {
	cb func()
}

func (t *MockCbTimer) WaitCb(d time.Duration, cb func()) error {

	if t.cb != nil {
		return TimerWaitedError
	}
	t.cb = cb
	return nil
}

func (t *MockCbTimer) Cancel() error {
	if t.cb == nil {
		return TimerNotWaited
	}
	t.cb = nil
	return nil
}

func (t *MockCbTimer) ForceDone() error {
	if t.cb == nil {
		return TimerNotWaited
	}

	// THIS IS NECESSARY TO POTENTIALLY RE-START THE TIMER FROM THIS VERY
	// CALLBACK.
	cb := t.cb
	t.cb = nil

	cb()
	return nil
}
