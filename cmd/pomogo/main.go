package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/FernandoAFS/pomogo/config"
	"os"
)

//go:generate sh -c "printf %s $(git rev-parse HEAD) > commit.txt"
//go:embed commit.txt
var Commit string

//go:generate sh -c "printf %s $(git describe) > version.txt"
//go:embed version.txt
var Version string

var helpMessage = "No command. Use `server`, `client` or `version`."

func onErr(err error) {
	if err == nil {
		return
	}
	fmt.Fprintln(os.Stderr, err)
	os.Exit(2)
}

func main() {
	if len(os.Args) <= 1 {
		fmt.Fprintln(os.Stderr, helpMessage)
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
	case "version":
		fmt.Printf(
			"Version: %s\nCommit: %s\n",
			Version,
			Commit,
		)
	default:
		fmt.Fprintf(
			os.Stderr,
			"Unknown command %s. %s\n",
			command,
			helpMessage,
		)
	}
}
