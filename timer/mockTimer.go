package timer

// Timer wrapper. Object meant for easier testing.
// This version of the timer keeps the done chanel and makes the cancel and
// forceDone methos blocking.

import (
	"sync"
	"time"
)

// TIMER WITH LOCK PROTECTION
type MockTimer struct {
	doneChan chan bool
	isWaited bool
	now      time.Time
	locker   sync.Mutex
}

func MockTimerFactory(refTime time.Time) *MockTimer {
	return &MockTimer{
		doneChan: make(chan bool),
		now:      refTime,
	}
}

func (t *MockTimer) setup() error {
	t.locker.Lock()
	defer t.locker.Unlock()

	// UNLIKE IN "TRUE" IMPLEMENTATION ONLY KEEP BOOLEAN. THIS IS DONE SO IT'S
	// SAFE TO WAIT FOR "FORCE DONE" BEFORE THE TIMER IS EVEN WAITED. BUT WE
	// CAN'T 2 LISTENERS.
	if t.isWaited {
		return TimerWaitedError
	}

	t.isWaited = true
	return nil
}

// THREAD SAFE WAY TO REMOVE CANCEL CHANEL.
func (t *MockTimer) teardown() error {
	t.locker.Lock()
	defer t.locker.Unlock()

	if !t.isWaited {
		return TimerWaitedError
	}

	t.isWaited = false
	return nil
}

func (t *MockTimer) cbWaitRoutine(then time.Time, cb func()) {
	done := <-t.doneChan
	// FIRST THING. CHANEL IS NO LONGER WAITED
	t.teardown()
	if !done {
		return
	}
	cb()
}

func (t *MockTimer) WaitCb(d time.Duration, cb func()) error {
	if err := t.setup(); err != nil {
		return err
	}
	then := t.now.Add(d)
	go t.cbWaitRoutine(then, cb)
	return nil
}

func (t *MockTimer) Cancel() error {
	// UNLIKE IN "TRUE" TIMER THIS BLOCKS THE ROUTINE
	t.doneChan <- false

	t.locker.Lock()
	defer t.locker.Unlock()
	t.isWaited = false
	return nil
}

func (t *MockTimer) ForceDone() error {
	// UNLIKE IN "TRUE" TIMER THIS BLOCKS THE ROUTINE
	t.doneChan <- true

	t.locker.Lock()
	defer t.locker.Unlock()
	t.isWaited = false
	return nil
}
