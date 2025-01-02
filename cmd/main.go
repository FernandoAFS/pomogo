package main

import (
	"encoding/json"
	"fmt"
	"os"
	"pomogo/config"
)

func onErr(err error) {
	if err == nil {
		return
	}
	fmt.Fprintln(os.Stderr, err)
	os.Exit(2)
}

func main() {
	if len(os.Args) <= 1 {
		fmt.Fprintln(os.Stderr, "No command. Use `server` or `client`.")
		os.Exit(2)
	}
	subArgs := os.Args[2:]
	command := os.Args[1]

	switch command {
	case "server":
		srvCfg, err := config.ServerCmdArgParse(subArgs...)
		onErr(err)
		onErr(srvCfg.HttpListen())
	case "client":
		clCfg, err := config.ClientCmdArgParse(subArgs...)
		onErr(err)
		st, err := clCfg.Run()
		onErr(err)
		// TODO: IMPROVE ON JSON PRINT STYLE...
		// fmt.Println(st)
		r, err := json.MarshalIndent(st, "", "\t")
		onErr(err)
		fmt.Println(string(r))
	default:
		fmt.Fprintf(
			os.Stderr,
			"Unknown command %s. Use `server` or `client`",
			command,
		)
	}
}
