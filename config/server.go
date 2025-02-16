package config

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/FernandoAFS/pomogo/controller"
	"github.com/FernandoAFS/pomogo/server"
	"github.com/FernandoAFS/pomogo/session"
	"github.com/FernandoAFS/pomogo/timer"
)

type ServerConfig struct {
	nSessions          int
	listenProto        string
	listenAddress      string
	workDuration       time.Duration
	shortBreakDuration time.Duration
	longBreakDuration  time.Duration
	command            string
}

func ServerCmdArgParse(args ...string) (*ServerConfig, error) {
	fs := flag.NewFlagSet("server", flag.ExitOnError)

	// PROBABLY IMPROVE ON ERROR MANAGEMENT...
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	nSessions := fs.Int(
		"work_sessions",
		4,
		"Number of break sessions.",
	)
	listenProto := fs.String(
		"protocol",
		"unix",
		"Protocol for communications. Use unix for file or tcp for tcp/ip.",
	)

	listenAddress := fs.String(
		"address",
		homeDir+"/.pomogo.socket",
		"Address for communications. Use unix for file or tcp for tcp/ip",
	)

	workDuration := fs.Duration(
		"work_duration",
		25*60_000000000, // 25m
		"Duration of work session.",
	)

	shortBreakDuration := fs.Duration(
		"short_break_duration",
		5*60_000000000, // 25m
		"Duration of short break.",
	)

	longBreakDuration := fs.Duration(
		"long_break_duration",
		15*60_000000000, // 25m
		"Duration of long break.",
	)

	command := fs.String(
		"event_command",
		"",
		"Command to be runned on every controller event (but error)",
	)

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	// VALIDATION LOGIC
	// TODO: MOVE VALIDATION TO ITS OWN METHOD.

	if *listenProto != "unix" && *listenProto != "tcp" {
		return nil, fmt.Errorf("invalid argument: %s", *listenProto)
	}
	// TODO: CROSS CHECK PROTOCOL AND ADDRESS.

	return &ServerConfig{
		nSessions:          *nSessions,
		listenProto:        *listenProto,
		listenAddress:      *listenAddress,
		workDuration:       *workDuration,
		shortBreakDuration: *shortBreakDuration,
		longBreakDuration:  *longBreakDuration,
		command:            *command,
	}, nil
}

func (sc *ServerConfig) sessionFactory() session.PomoSessionIface {
	return &session.PomoSession{
		WorkSessionsBreak: sc.nSessions,
	}
}

func (sc *ServerConfig) timerFactory() timer.PomoTimerIface {
	return new(timer.PomoTimer)
}

func (sc *ServerConfig) durationFactory() session.SessionStateDurationFactory {
	return session.DurationFactory(
		sc.workDuration,
		sc.shortBreakDuration,
		sc.longBreakDuration,
	)
}

func (sc *ServerConfig) controllerFactory() (controller.PomoControllerIface, error) {

	options := []controller.PomoControllerOption{
		controller.PomoControllerSessionOpt(sc.sessionFactory),
		controller.PomoControllerTimerOpt(sc.timerFactory),
		controller.PomoControllerDurationF(sc.durationFactory),
	}

	if sc.command != "" {
		options = append(options, controller.PomoControllerHook(sc.command))
	}

	return controller.ControllerFactory(
		options...,
	)
}

func (sc *ServerConfig) controllerFactoryPanic() controller.PomoControllerIface {
	ctrl, err := sc.controllerFactory()
	if err != nil {
		panic(err)
	}
	return ctrl
}

func (sc *ServerConfig) containerFactory() *controller.SingleControllerContainer {
	return &controller.SingleControllerContainer{
		ControllerFactory: sc.controllerFactoryPanic,
	}
}

func (sc *ServerConfig) serverFactory() (*server.SingleSessionServer, error) {
	return server.SingleSessionServerFactory(
		server.SingleServerContainerOpt(sc.containerFactory),
	)
}

func (sc *ServerConfig) runServerCtx() server.SServerFuncOpt {
	if sc.listenProto == "unix" {
		return server.SingleServerRpcUnixRegOpt(
			sc.listenAddress,
			rpc.NewServer,
			func(l net.Listener, s *rpc.Server) error {

				exit := make(chan os.Signal, 1)
				signal.Notify(exit, os.Interrupt, syscall.SIGTERM)

				// TODO: build a more elegant solution...
				go func() {
					<-exit
					err := l.Close()
					if err != nil {
						fmt.Println(err)
					}
				}()

				return http.Serve(l, s)
			},
		)
	}

	return server.SingleServerRpcRegisterOpt(
		sc.listenProto,
		sc.listenAddress,
		rpc.NewServer,
		func(l net.Listener, s *rpc.Server) error {
			return http.Serve(l, s)
		},
	)
}

// Run appropiate server through http synchronously
func (sc *ServerConfig) HttpListen() error {
	run_srv := sc.runServerCtx()
	srv, err := sc.serverFactory()
	if err != nil {
		return err
	}

	undo, err := run_srv(srv)
	if err != nil {
		return err
	}

	_, err = undo(srv)
	if err != nil {
		return err
	}

	return nil
}
