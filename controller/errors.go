package controller

import "errors"

var StoppedTimer = errors.New("Cannot execute action on stopped timer")
var PausedTimer = errors.New("Cannot execute action on paused timer")
var RunningTimer = errors.New("Cannot execute action on running timer")
