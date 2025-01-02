package controller

import "testing"

// ========
// FIXTURES
// ========

func getTestContainer(ctrl PomoControllerIface) SingleControllerContainer {
	return SingleControllerContainer{
		ControllerFactory: func() PomoControllerIface {
			return ctrl
		},
	}
}

// =====
// TESTS
// =====

// Standard case. Create, fetch and remove
func TestControllerContainerUse(t *testing.T) {
	ctrl := PomoController{}
	container := getTestContainer(&ctrl)

	container.CreateController()

	ctrlRes := container.GetController()

	if &ctrl != ctrlRes {
		t.Fatalf("Controller reference mismatch")
	}

	if err := container.RemoveController(); err != nil {
		t.Fatal(err)
	}

}

func TestControllerContainerEarlyRemove(t *testing.T) {
	// Check error on early remove

	ctrl := PomoController{}
	container := getTestContainer(&ctrl)

	err := container.RemoveController()
	if err != ErrNoControllerError {
		t.Fatal(err)
	}
}
