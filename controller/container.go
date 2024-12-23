// THIS CONTAINS THE LOGIC OF A SINGLE SESSION MANAGER WRAPPED UNDER THE
// CONTROLLER INTERFACE

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
func (c *SingleControllerContainer) CreateController() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.controller != nil {
		return ExistintgControllerError
	}

	c.controller = c.ControllerFactory()
	return nil
}

// Return existing controller. Return error if none exist
func (c *SingleControllerContainer) GetController() (PomoControllerIface, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if c.controller == nil {
		return nil, NoControllerError
	}

	return c.controller, nil
}

// Remove reference to instance. Return error if none exist.
func (c *SingleControllerContainer) RemoveController() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.controller == nil {
		return NoControllerError
	}

	c.controller = nil
	return nil
}
