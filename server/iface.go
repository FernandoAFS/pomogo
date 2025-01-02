package server

import (
	pomoController "pomogo/controller"
)

type pomoStatus = pomoController.PomoControllerStatus
type pomoCtrl = pomoController.PomoControllerIface

type PomogoSessionServer interface {
	Status(
		request struct{},
		reply *pomoController.PomoControllerStatus,
	) error
	Pause(
		request struct{},
		reply *pomoController.PomoControllerStatus,
	) error
	Play(
		request struct{},
		reply *pomoController.PomoControllerStatus,
	) error
	Skip(
		request struct{},
		reply *pomoController.PomoControllerStatus,
	) error
	Stop(
		request struct{},
		reply *pomoController.PomoControllerStatus,
	) error
}

type PomogoClient interface {
	Status() (*pomoStatus, error)
	Pause() (*pomoStatus, error)
	Play() (*pomoStatus, error)
	Skip() (*pomoStatus, error)
	Stop() (*pomoStatus, error)
}
