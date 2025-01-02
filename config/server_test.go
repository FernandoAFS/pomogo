package config

import "testing"

// Extremely basic test. Controlled inputs lead to no error
func TestServerConfigValidation(t *testing.T) {

	sc := ServerConfig{
		nSessions:          4,
		listenProto:        "unix",
		listenAddress:      "/tmp/pomogo.socket",
		workDuration:       25 * 60_000000000,
		shortBreakDuration: 5 * 60_000000000,
		longBreakDuration:  15 * 60_000000000,
	}

	if _, err := sc.controllerFactory(); err != nil {
		t.Fatal(err)
	}

	if _, err := sc.serverFactory(); err != nil {
		t.Fatal(err)
	}
}
