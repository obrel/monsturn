package main

import (
	"os"

	"github.com/obrel/monsturn/cmd"
)

func main() {
	cmd.Init()

	if err := cmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
