// Common logic on what to do on error
// TODO: IMPROVE ON ERROR MANAGEMENT. Hooks should return an error. This code should not output to stderr.

package controller

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

// ======================
// COMMAND EXECUTOR LOGIC
// ======================

func onError(err error) {
	if err == nil {
		return
	}
	fmt.Fprint(os.Stderr, err)
}

func genCommand(command string, at time.Time, status string, eventType string) *exec.Cmd {
	cmd := exec.Command(command)
	cmd.Env = append(
		os.Environ(),
		"POMOGO_AT="+at.String(),
		"POMOGO_STATUS="+status,
		"POMO_EVENT="+eventType,
	)
	return cmd
}

func PlayExecHook(command string) func(event PomoControllerEventArgsPlay) {
	return func(event PomoControllerEventArgsPlay) {
		cmd := genCommand(
			command,
			event.At,
			event.CurrentState.String(),
			"Play",
		)
		go onError(cmd.Run())
	}
}

func StopExecHook(command string) func(event PomoControllerEventArgsStop) {
	return func(event PomoControllerEventArgsStop) {
		cmd := genCommand(
			command,
			event.At,
			event.CurrentState.String(),
			"Stop",
		)
		go onError(cmd.Run())
	}
}

func PauseExecHook(command string) func(event PomoControllerEventArgsPause) {
	return func(event PomoControllerEventArgsPause) {
		cmd := genCommand(
			command,
			event.At,
			event.CurrentState.String(),
			"Pause",
		)
		go onError(cmd.Run())
	}
}

func NextStateExecHook(command string) func(event PomoControllerEventArgsNextState) {
	return func(event PomoControllerEventArgsNextState) {
		cmd := genCommand(
			command,
			event.At,
			event.NextState.String(),
			"EndOfState",
		)
		go onError(cmd.Run())
	}
}

func ErrorExecHook(command string) func(event error) {
	return func(event error) {
		cmd := genCommand(
			command,
			time.Now(),
			event.Error(),
			"Error",
		)
		go onError(cmd.Run())
	}
}
