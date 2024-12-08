package timer

import (
	"sync"
	"time"
)

// TIMER WITH LOCK PROTECTION
type PomoTimer struct {
	cancelChan chan struct{}
	locker     sync.Locker
}

// THREAD SAFE WAY TO CREATE A NEW CANCEL CHANEL.
func (t *PomoTimer) setup() error {
	t.locker.Lock()
	defer t.locker.Unlock()

	if t.cancelChan != nil {
		return TimerWaitedError
	}

	t.cancelChan = make(chan struct{})
	return nil
}

// THREAD SAFE WAY TO REMOVE CANCEL CHANEL.
func (t *PomoTimer) teardown() bool {
	t.locker.Lock()
	defer t.locker.Unlock()

	if t.cancelChan == nil {
		return false
	}

	// GRACEFUL DELETION
	close(t.cancelChan)
	return true
}

func (t *PomoTimer) cbWaitRoutine(d time.Duration, cb func()) {
	select {
	case <-time.After(d):
		t.teardown()
		cb() // ASYNC EXECUTION OF THE CALLBACK AND INMEDIATE TEARDOWN.
	case <-t.cancelChan:
		t.teardown()
	}
}

func (t *PomoTimer) WaitCb(d time.Duration, cb func()) error {
	// CANNOT WAIT TWICE.
	if err := t.setup(); err != nil {
		return err
	}
	go t.cbWaitRoutine(d, cb)
	return nil
}

func (t *PomoTimer) Cancel(d time.Duration) error {
	t.locker.Lock()
	defer t.locker.Unlock()

	if t.cancelChan == nil {
		return TimerNotWaited
	}

	t.cancelChan <- struct{}{}
	t.cancelChan = nil
	return nil
}
