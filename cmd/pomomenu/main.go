package main

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

//go:embed script.sh
var PomomenuScript string

func main() {
	// pass parameters to your script as a safe way.
	c := exec.Command("sh", "-s", "-", "-la", "/etc")

	// use $1, $2, ... $@ as usual
	c.Stdin = strings.NewReader(PomomenuScript)

	b, err := c.Output()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println(string(b))
}
