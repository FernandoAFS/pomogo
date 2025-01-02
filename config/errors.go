package config

import (
	"errors"
	"fmt"
)

var ErrInvalidArg = errors.New("invalid argument: ")

// TODO: THIS IS A BAD IDEA. USE FMT.ERRORF INSTEAD
func NewInvalidArgError(arg ...interface{}) error {
	argMsg := fmt.Sprint(arg...)
	return errors.Join(ErrInvalidArg, errors.New(argMsg))
}
