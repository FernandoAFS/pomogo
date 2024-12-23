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

func TestControllerContainerUse(t *testing.T) {
	// Standard case. Create, fetch and remove
	ctrl := PomoController{}
	container := getTestContainer(&ctrl)

	if err := container.CreateController(); err != nil {
		t.Fatal(err)
	}

	ctrlRes, err := container.GetController()

	if err != nil {
		t.Fatal(err)
	}

	if &ctrl != ctrlRes {
		t.Fatalf("Controller reference mismatch")
	}

	if err := container.RemoveController(); err != nil {
		t.Fatal(err)
	}

}

func TestControllerContainerEarlyGet(t *testing.T) {
	// Check error on early return

	ctrl := PomoController{}
	container := getTestContainer(&ctrl)

	_, err := container.GetController()
	if err != NoControllerError {
		t.Fatal(err)
	}
}

func TestControllerContainerEarlyRemove(t *testing.T) {
	// Check error on early remove

	ctrl := PomoController{}
	container := getTestContainer(&ctrl)

	err := container.RemoveController()
	if err != NoControllerError {
		t.Fatal(err)
	}
}

func TestControllerContainerDoubleCreate(t *testing.T) {
	// Check error on create over existing container

	ctrl := PomoController{}
	container := getTestContainer(&ctrl)

	if err := container.CreateController(); err != nil {
		t.Fatal(err)
	}

	if err := container.CreateController(); err != ExistintgControllerError {
		t.Fatal(err)
	}
}
