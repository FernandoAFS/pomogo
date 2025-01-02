package cli

import (
	"flag"
	"os"
	"pomogo/controller"
	"pomogo/server"
	"strings"
)

type ClientConfig struct {
	connectProto   string
	connectAddress string
	action         string
}

// Generate object from flags.
func ClientCmdArgParse(args ...string) (*ClientConfig, error) {
	fs := flag.NewFlagSet("client", flag.ExitOnError)

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	connectProto := fs.String(
		"protocol",
		"unix",
		"Protocol for communications. Use unix for file or tcp for tcp/ip.",
	)

	listenAddressDef := homeDir + "/pomogo.socket"
	connectAddress := fs.String(
		"address",
		listenAddressDef,
		"Address for communications. Use unix for file or tcp for tcp/ip.",
	)

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	action := fs.Arg(0)

	cc := &ClientConfig{
		connectProto:   *connectProto,
		connectAddress: *connectAddress,
		action:         action,
	}

	return cc, nil
}

// Perform the action and return status or error
func (cc *ClientConfig) Run() (*controller.PomoControllerStatus, error) {
	cl, err := server.PomogoRpcClientFactory(cc.connectProto, cc.connectAddress)

	if err != nil {
		return nil, err
	}

	switch strings.ToLower(cc.action) {
	case "status":
		return cl.Status()
	case "pause":
		return cl.Pause()
	case "play":
		return cl.Play()
	case "skip":
		return cl.Skip()
	case "stop":
		return cl.Stop()
	}

	return nil, NewInvalidArgError(cc.action)
}
