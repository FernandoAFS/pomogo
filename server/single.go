package server

import (
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
type SingleSessionServer struct {
	container *pomoController.SingleControllerContainer
}

// 100% PRIVATE DRY METHOD
func (c *SingleSessionServer) doNowCb(
	cb func(now time.Time, ctrl pomoCtrl) error,
) (*pomoController.PomoControllerStatus, error) {

	ctrl, err := c.container.GetController()

	if err != nil {
		return nil, err
	}

	if err := cb(time.Now(), ctrl); err != nil {
		return nil, err
	}

	status := ctrl.Status()
	return &status, nil
}

func (c *SingleSessionServer) Status(
	request struct{},
	reply *pomoController.PomoControllerStatus,
) error {
	ctrl, err := c.container.GetController()
	if err != nil {
		return err
	}
	status := ctrl.Status()
	reply = &status
	return nil
}

func (c *SingleSessionServer) Pause(
	request struct{},
	reply *pomoController.PomoControllerStatus,
) error {

	reply, err := c.doNowCb(func(now time.Time, ctrl pomoCtrl) error {
		return ctrl.Pause(now)
	})
	return err
}

// CREATE NEW INSTANCE AND START CONTROLLER COUNTING.
func (c *SingleSessionServer) Play(
	request struct{},
	reply *pomoController.PomoControllerStatus,
) error {
	if err := c.container.CreateController(); err != nil {
		return err
	}

	reply, err := c.doNowCb(func(now time.Time, ctrl pomoCtrl) error {
		return ctrl.Play(now)
	})
	return err
}

func (c *SingleSessionServer) Skip(
	request struct{},
	reply *pomoController.PomoControllerStatus,
) error {
	reply, err := c.doNowCb(func(now time.Time, ctrl pomoCtrl) error {
		return ctrl.Skip(now)
	})
	return err
}

func (c *SingleSessionServer) Stop(
	request struct{},
	reply *pomoController.PomoControllerStatus,
) error {
	reply, err := c.doNowCb(func(now time.Time, ctrl pomoCtrl) error {
		return ctrl.Stop(now)
	})
	return err
}

// GIVEN A SERVER START LISTENING LISTENING SYNCHRONOUSLY
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
	var resp *pomoStatus
	callName := DefaultServerName + "." + method

	if err := c.client.Call(callName, struct{}{}, &resp); err != nil {
		return nil, err
	}

	return resp, nil
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
