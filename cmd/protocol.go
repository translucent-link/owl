package cmd

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	model "github.com/translucent-link/owl/graph/model"

	"github.com/urfave/cli/v2"
)

func protocolRegisterHandler(c *cli.Context) error {
	db, err := model.DbConnect()
	if err != nil {
		return err
	}
	defer db.Close()
	store := model.NewProtocolStore(db)

	abiBytes, err := os.ReadFile(c.String("abiFilepath"))
	if err != nil {
		return errors.Wrapf(err, "Reading ABI file: %s", c.String("abiFilepath"))
	}

	newProtocol, err := store.CreateProtocol(model.NewProtocol{
		Name: c.String("name"),
		Abi:  string(abiBytes),
	})
	if err != nil {
		return err
	}
	fmt.Printf("Registered new protocol: %s\n", newProtocol.Name)
	return nil
}

func protocolDeleteHandler(c *cli.Context) error {
	db, err := model.DbConnect()
	if err != nil {
		return err
	}
	defer db.Close()
	store := model.NewProtocolStore(db)
	err = store.DeleteByName(c.String("name"))
	if err != nil {
		return err
	}
	fmt.Printf("Deleted protocol: %s\n", c.String("name"))
	return nil
}

func protocolListHandler(c *cli.Context) error {
	db, err := model.DbConnect()
	if err != nil {
		return err
	}
	defer db.Close()
	store := model.NewProtocolStore(db)
	protocols, err := store.All()
	if err != nil {
		return err
	}
	for _, protocol := range protocols {
		fmt.Printf("%s\n", protocol.Name)
	}
	return nil
}

var registerProtocolCmd = &cli.Command{
	Name:   "register",
	Usage:  "registers a new protocol",
	Action: protocolRegisterHandler,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "name",
			Usage:    "e.g Compound",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "abiFilepath",
			Usage:    "e.g ../abi/erc20.json",
			Required: true,
		},
	},
}

var deleteProtocolCmd = &cli.Command{
	Name:   "delete",
	Usage:  "deletes a protocol and all associated protocol instances and events",
	Action: protocolDeleteHandler,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "name",
			Usage:    "e.g polygon",
			Required: true,
		},
	},
}

var listProtocolCmd = &cli.Command{
	Name:   "ls",
	Usage:  "lists all protocols",
	Action: protocolListHandler,
}

var ProtocolCommand = &cli.Command{
	Name:  "protocol",
	Usage: "registers and manages protocols and their ABIs",
	Flags: []cli.Flag{},
	Subcommands: []*cli.Command{

		registerProtocolCmd, deleteProtocolCmd, listProtocolCmd,
	},
}
