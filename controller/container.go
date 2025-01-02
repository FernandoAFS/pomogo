// This contains the logic of a single session manager wrapped under the
// controller interface

package controller

import "sync"

// Single reference instance of controller container.
// Creates instance of controller on first request or after delete.
// Typically you should call Stop() method on controller to avoid having goroutines open.
type SingleControllerContainer struct {
	ControllerFactory func() PomoControllerIface
	controller        PomoControllerIface
	mutex             sync.RWMutex
}

// Create new controlle instance. Return error if one already exists
func (c *SingleControllerContainer) CreateController() PomoControllerIface {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.controller == nil {
		c.controller = c.ControllerFactory()
	}

	return c.controller
}

// Return existing controller. Return error if none exist
func (c *SingleControllerContainer) GetController() PomoControllerIface {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.controller
}

// Remove reference to instance. Return error if none exist.
func (c *SingleControllerContainer) RemoveController() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.controller == nil {
		return ErrNoControllerError
	}

	c.controller = nil
	return nil
}
