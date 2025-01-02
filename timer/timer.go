package timer

import (
	"sync"
	"time"
)

// TIMER WITH LOCK PROTECTION
type PomoTimer struct {
	cancelChan chan struct{}
	locker     sync.Mutex
}

// THREAD SAFE WAY TO CREATE A NEW CANCEL CHANEL.
func (t *PomoTimer) setup() error {
	t.locker.Lock()
	defer t.locker.Unlock()

	if t.cancelChan != nil {
		return ErrTimerWaited
	}

	t.cancelChan = make(chan struct{})
	return nil
}

func (t *PomoTimer) cbWaitRoutine(d time.Duration, cb func()) {

	teardown := func() bool {
		t.locker.Lock()
		defer t.locker.Unlock()
		if t.cancelChan == nil {
			return false
		}
		close(t.cancelChan)
		t.cancelChan = nil
		return true
	}

	select {
	case <-time.After(d):
		// If chan was terminated between timer event and lock capture skip
		if !teardown() {
			return
		}
		cb()
		return
	case <-t.cancelChan:
		// PROBABLY LOG THIS.
		return
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

func (t *PomoTimer) Cancel() error {
	t.locker.Lock()
	defer t.locker.Unlock()

	if t.cancelChan == nil {
		return ErrTimerNotWaited
	}

	t.cancelChan <- struct{}{}
	close(t.cancelChan)
	t.cancelChan = nil
	return nil
}
