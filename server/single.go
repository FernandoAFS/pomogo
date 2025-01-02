package server

import (
	"log/slog"
	"net"
	"net/http"
	"net/rpc"
	pomoController "pomogo/controller"
	"time"
)

// ======
// SERVER
// ======

// LIFECYCLE MANAGEMENT OF CONTROLLER.
// It may be better to skip the container step and manage the lifecycle
// directly from the server object.
type SingleSessionServer struct {
	container *pomoController.SingleControllerContainer
}

// 100% private dry method
func (c *SingleSessionServer) doNowCb(
	cb func(ctrl pomoCtrl) error,
) error {

	ctrl := c.container.GetController()
	if ctrl == nil {
		return pomoController.NoControllerError
	}

	if err := cb(ctrl); err != nil {
		return err
	}
	return nil
}

func (c *SingleSessionServer) Status(
	request struct{},
	reply *pomoController.PomoControllerStatus,
) error {
	ctrl := c.container.GetController()
	if ctrl == nil {
		return pomoController.NoControllerError
	}
	*reply = ctrl.Status()
	return nil
}

func (c *SingleSessionServer) Pause(
	request struct{},
	reply *pomoController.PomoControllerStatus,
) error {
	now := time.Now()
	return c.doNowCb(
		func(ctrl pomoCtrl) error {
			if err := ctrl.Pause(now); err != nil {
				return err
			}
			*reply = ctrl.Status()
			return nil
		})
}

// CREATE NEW INSTANCE AND START CONTROLLER COUNTING.
func (c *SingleSessionServer) Play(
	request struct{},
	reply *pomoController.PomoControllerStatus,
) error {
	ctrl := c.container.CreateController()
	now := time.Now()
	if err := ctrl.Play(now); err != nil {
		return err
	}
	*reply = ctrl.Status()
	return nil
}

func (c *SingleSessionServer) Skip(
	request struct{},
	reply *pomoController.PomoControllerStatus,
) error {
	now := time.Now()
	return c.doNowCb(
		func(ctrl pomoCtrl) error {
			if err := ctrl.Skip(now); err != nil {
				return err
			}
			*reply = ctrl.Status()
			return nil
		})
}

func (c *SingleSessionServer) Stop(
	request struct{},
	reply *pomoController.PomoControllerStatus,
) error {
	now := time.Now()
	return c.doNowCb(
		func(ctrl pomoCtrl) error {
			if err := ctrl.Stop(now); err != nil {
				return err
			}
			*reply = ctrl.Status()
			return nil
		})
}

// Given a server start listening listening synchronously
func SingleSessionServerStart(protocol, address string, wrapper *SingleSessionServer) error {
	// Name to be registered.
	rpc.RegisterName(DefaultServerName, wrapper)
	rpc.HandleHTTP()
	l, err := net.Listen(protocol, address)
	if err != nil {
		return err
	}
	// COMMUNICATE WITH TEST CASE.
	http.Serve(l, nil)
	return nil
}

// ======
// CLIENT
// ======

// Direct wrapper around the client. The lifecycle of the object should not be
// longer than the connection:
type SingleSessionClient struct {
	client *rpc.Client
}

// Simply call a method given the string name and return the response as a
// pomodoro status
func (c *SingleSessionClient) callMethod(method string) (*pomoStatus, error) {
	var resp pomoStatus
	callName := DefaultServerName + "." + method

	slog.Debug("Making request", "method", callName, "response", resp)

	if err := c.client.Call(callName, struct{}{}, &resp); err != nil {
		return nil, err
	}

	slog.Debug("Successfull response", "status", resp)

	return &resp, nil
}

func (c *SingleSessionClient) Status() (*pomoStatus, error) {
	return c.callMethod("Status")
}

func (c *SingleSessionClient) Pause() (*pomoStatus, error) {
	return c.callMethod("Pause")
}

func (c *SingleSessionClient) Play() (*pomoStatus, error) {
	return c.callMethod("Play")
}

func (c *SingleSessionClient) Skip() (*pomoStatus, error) {
	return c.callMethod("Skip")
}

func (c *SingleSessionClient) Stop() (*pomoStatus, error) {
	return c.callMethod("Stop")
}

// Initializes rpc client and returns wrapper
func PomogoRpcClientFactory(protocol, address string) (*SingleSessionClient, error) {

	client, err := rpc.DialHTTP(protocol, address)
	if err != nil {
		return nil, err
	}

	cl := SingleSessionClient{
		client: client,
	}

	return &cl, nil
}
