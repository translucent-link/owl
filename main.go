package main

import (
	"log"
	"os"

	"github.com/translucent-link/owl/cmd"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			cmd.BlkCommand,
			cmd.ListenCommand,
			cmd.ScanCommand,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
