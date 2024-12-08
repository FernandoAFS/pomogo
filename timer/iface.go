package timer

import (
	"time"
)

type PomoTimerIface interface {
	WaitCb(d time.Duration, cb func()) error
	Cancel() error
}
