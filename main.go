package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/translucent-link/owl/cmd"
	"github.com/urfave/cli/v2"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func setupDatabase(dbURL string) {
	m, err := migrate.New("file://./db/migrations", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil && err.Error() != "no change" {
		log.Fatal(err)
	}
}

func main() {
	_ = godotenv.Load()

	dbURL := os.Getenv("DATABASE_URL")
	if len(dbURL) == 0 {
		log.Fatal("Please specify DATABASE_URL environment variable.")
	}
	setupDatabase(dbURL)

	app := &cli.App{
		Commands: []*cli.Command{
			cmd.BlkCommand,
			cmd.ListenCommand,
			cmd.ScanCommand,
			cmd.ServerCommand,
			cmd.AbiCommand,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
