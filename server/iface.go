package server

import (
	pomoController "pomogo/controller"
)

type pomoStatus = pomoController.PomoControllerStatus
type pomoCtrl = pomoController.PomoControllerIface

type PomogoClient interface {
	Status() (*pomoStatus, error)
	Pause() (*pomoStatus, error)
	Play() (*pomoStatus, error)
	Skip() (*pomoStatus, error)
	Stop() (*pomoStatus, error)
}
