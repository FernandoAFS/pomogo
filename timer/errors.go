package timer

import "errors"

var TimerWaitedError = errors.New("Timer aleready waited. Cannot wait twice.")
var TimerNotWaited = errors.New("Timer not waited. Cannot cancel.")
