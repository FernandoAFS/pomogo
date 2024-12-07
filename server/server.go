package rpc

import (
	pomoController "pomogo/controller"
	"time"
)

// LIFECYCLE MANAGEMENT OF CONTROLLER.
type SingleSessionContainer struct {
	controllerFactory func() pomoController.PomoControllerIface
	controller        pomoController.PomoControllerIface
}

// 100% PRIVATE DRY METHOD
func (c *SingleSessionContainer) doNowCb(
	cb func(now time.Time) error,
) (*pomoController.PomoControllerStatus, error) {

	if c.controller == nil {
		return nil, NoControllerError
	}

	if err := cb(time.Now()); err != nil {
		return nil, err
	}

	status := c.controller.Status()
	return &status, nil
}

func (c *SingleSessionContainer) Status(
	request struct{},
	reply *pomoController.PomoControllerStatus,
) error {

	if c.controller == nil {
		return NoControllerError
	}

	status := c.controller.Status()
	reply = &status
	return nil
}

func (c *SingleSessionContainer) Pause(
	request struct{},
	reply *pomoController.PomoControllerStatus,
) error {

	st, err := c.doNowCb(func(now time.Time) error {
		return c.controller.Pause(now)
	})

	reply = st
	return err
}

func (c *SingleSessionContainer) Play(
	request struct{},
	reply *pomoController.PomoControllerStatus,
) error {

	// CREATE CONTROLLER ON FIRST PLAY.
	if c.controller == nil {
		c.controller = c.controllerFactory()
	}

	st, err := c.doNowCb(func(now time.Time) error {
		return c.controller.Play(now)
	})
	reply = st
	return err
}

func (c *SingleSessionContainer) Skip(
	request struct{},
	reply *pomoController.PomoControllerStatus,
) error {

	st, err := c.doNowCb(func(now time.Time) error {
		return c.controller.Skip(now)
	})
	reply = st
	return err
}

func (c *SingleSessionContainer) Stop(
	request struct{},
	reply *pomoController.PomoControllerStatus,
) error {

	st, err := c.doNowCb(func(now time.Time) error {
		return c.controller.Stop(now)
	})
	reply = st

	// REMOVE REFERENCE TO CONTROLLER REGARDLESS OF ERROR STATUS.
	c.controller = nil
	return err
}
