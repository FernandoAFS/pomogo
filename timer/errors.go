package timer

import "errors"

var ErrTimerWaited = errors.New("timer aleready waited, cannot wait twice")
var ErrTimerNotWaited = errors.New("timer not waited, cannot cancel")
