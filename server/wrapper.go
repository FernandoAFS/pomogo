package server

// Very crude servrer wrapper. Should be improved upon by including contexts
// and token to each request

import (
	pomoController "github.com/FernandoAFS/pomogo/controller"
	"log/slog"
)

type SessionWrapper struct {
	serverSession PomogoSessionServer
}

func (sw *SessionWrapper) Status(
	request struct{},
	reply *pomoController.PomoControllerStatus,
) error {
	slog.Info("Status Request")
	err := sw.serverSession.Status(request, reply)
	slog.Info("Status Response", "reply", reply, "err", err)
	return err
}

func (sw *SessionWrapper) Pause(
	request struct{},
	reply *pomoController.PomoControllerStatus,
) error {
	slog.Info("Status Request")
	err := sw.serverSession.Pause(request, reply)
	slog.Info("Status Response", "reply", reply, "err", err)
	return err
}

func (sw *SessionWrapper) Play(
	request struct{},
	reply *pomoController.PomoControllerStatus,
) error {
	slog.Info("Status Request")
	err := sw.serverSession.Pause(request, reply)
	slog.Info("Status Response", "reply", reply, "err", err)
	return err
}

func (sw *SessionWrapper) Skip(
	request struct{},
	reply *pomoController.PomoControllerStatus,
) error {
	slog.Info("Status Request")
	err := sw.serverSession.Pause(request, reply)
	slog.Info("Status Response", "reply", reply, "err", err)
	return err
}

func (sw *SessionWrapper) Stop(
	request struct{},
	reply *pomoController.PomoControllerStatus,
) error {
	slog.Info("Status Request")
	err := sw.serverSession.Pause(request, reply)
	slog.Info("Status Response", "reply", reply, "err", err)
	return err
}
