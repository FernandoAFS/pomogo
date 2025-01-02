package controller

import "errors"

var ErrStoppedTimer = errors.New("cannot execute action on stopped timer")
var ErrPausedTimer = errors.New("cannot execute action on paused timer")
var ErrRunningTimer = errors.New("cannot execute action on running timer")
var ErrNoControllerError = errors.New("must create a controller first")
var ErrExistintgControllerError = errors.New("must remove existing controller")
