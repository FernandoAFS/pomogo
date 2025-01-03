package server

import (
	"errors"
	"net"
	"net/rpc"
	"os"
	pomoController "pomogo/controller"
)

const (
	DefaultServerName = "SingleSessionServer"
)

// =======
// FACTORY
// =======

type SServerFuncOpt func(*SingleSessionServer) (SServerFuncOpt, error)
type SClientFuncOpt func(*SingleSessionClient) (SClientFuncOpt, error)

// Creates a new *SingleSessionServer reference and runs it through every option.
func SingleSessionServerFactory(options ...SServerFuncOpt) (*SingleSessionServer, error) {
	serv := new(SingleSessionServer)
	for _, opt := range options {
		_, err := opt(serv)
		if err != nil {
			return nil, err
		}
	}
	return serv, nil
}

// Create single session rpc client and run it through every option
func SingleSessionClientFactory(options ...SClientFuncOpt) (*SingleSessionClient, error) {
	client := new(SingleSessionClient)
	for _, opt := range options {
		_, err := opt(client)
		if err != nil {
			return nil, err
		}
	}
	return client, nil
}

// --------------
// SERVER OPTIONS
// --------------

// Set controller container given a factory function.
func SingleServerContainerOpt(
	factory func() *pomoController.SingleControllerContainer,
) SServerFuncOpt {
	return func(ss *SingleSessionServer) (SServerFuncOpt, error) {
		prev := ss.container
		ss.container = factory()
		return func(ss *SingleSessionServer) (SServerFuncOpt, error) {
			ss.container = prev
			return SingleServerContainerOpt(factory), nil
		}, nil
	}
}

// Register the ssServer on an rpc server from server factory (can use
// rpc.Server directly as factory)
func SingleServerRpcRegisterOpt(
	protocol, address string,
	serverFactory func() *rpc.Server,
	onListen func(l net.Listener, s *rpc.Server) error,
) SServerFuncOpt {
	return func(ss *SingleSessionServer) (SServerFuncOpt, error) {
		server := serverFactory()
		if err := server.RegisterName(DefaultServerName, ss); err != nil {
			return nil, err
		}
		l, err := net.Listen(protocol, address)
		if err != nil {
			return nil, err
		}
		if err := onListen(l, server); err != nil {
			return nil, err
		}
		return func(ss *SingleSessionServer) (SServerFuncOpt, error) {
			if err := l.Close(); err != nil {
				return nil, err
			}
			return SingleServerRpcRegisterOpt(protocol, address, serverFactory, onListen), nil
		}, nil
	}
}

// Same as before but removes unix socket when done and returns error if socket already exists.
func SingleServerRpcUnixRegOpt(
	address string,
	serverFactory func() *rpc.Server,
	onListen func(l net.Listener, s *rpc.Server) error,
) SServerFuncOpt {
	reg := SingleServerRpcRegisterOpt("unix", address, serverFactory, onListen)
	return func(ss *SingleSessionServer) (SServerFuncOpt, error) {

		// RAISE ERROR IF THE SOCKET ALREADY EXISTS. PROBABLY BAD EXIT...
		if _, err := os.Stat(address); !errors.Is(err, os.ErrNotExist) {
			// CONSIDER RETURNING A CUSTOM ERROR...
			return nil, os.ErrExist
		}

		unreg, err := reg(ss)
		if err != nil {
			return nil, err
		}
		return func(ss *SingleSessionServer) (SServerFuncOpt, error) {
			if err := os.Remove(address); err != nil {
				return nil, err
			}
			return unreg(ss)
		}, nil
	}
}

// --------------
// CLIENT OPTIONS
// --------------

// Connect to http-rpc server.
func SingleClientRpcHttpConnect(protocol, address string) SClientFuncOpt {
	return func(cl *SingleSessionClient) (SClientFuncOpt, error) {
		client, err := rpc.DialHTTP(protocol, address)
		if err != nil {
			return nil, err
		}
		cl.client = client
		return func(cl *SingleSessionClient) (SClientFuncOpt, error) {
			if err := client.Close(); err != nil {
				return nil, err
			}
			return SingleClientRpcHttpConnect(protocol, address), nil
		}, nil
	}
}
