package cli

import (
	"errors"
	"fmt"
)

var InvalidArgError = errors.New("Invalid argument: ")

func NewInvalidArgError(arg ...interface{}) error {
	argMsg := fmt.Sprint(arg...)
	return errors.Join(InvalidArgError, errors.New(argMsg))
}
