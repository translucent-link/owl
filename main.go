package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/translucent-link/owl/cmd"
	"github.com/urfave/cli/v2"
)

func main() {
	_ = godotenv.Load()

	// setupDatabase()

	app := &cli.App{
		Commands: []*cli.Command{
			cmd.BlkCommand,
			cmd.ScanCommand,
			cmd.ServerCommand,
			cmd.AbiCommand,
			cmd.ChainCommand,
			cmd.ProtocolCommand,
			cmd.ProtocolInstanceCommand,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
