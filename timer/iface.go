package timer

import (
	"time"
)

type PomoTimerIface interface {
	WaitCb(d time.Duration, cb func(then time.Time)) error
	Cancel() error
}
