package cmd

import (
	"fmt"

	model "github.com/translucent-link/owl/graph/model"

	"github.com/urfave/cli/v2"
)

func protocolInstanceRegisterHandler(c *cli.Context) error {
	db, err := model.DbConnect()
	if err != nil {
		return err
	}
	defer db.Close()
	store := model.NewProtocolInstanceStore(db)

	newInstance, err := store.CreateProtocolInstance(model.NewProtocolInstance{
		Protocol:         c.String("protocol"),
		Chain:            c.String("chain"),
		ContractAddress:  c.String("contractAddress"),
		FirstBlockToRead: c.Int("firstBlockToRead"),
	})
	if err != nil {
		return err
	}
	fmt.Printf("Registered new instance at %s on %s\n", newInstance.ContractAddress, newInstance.Chain.Name)
	return nil
}

func protocolInstanceDeleteHandler(c *cli.Context) error {
	db, err := model.DbConnect()
	if err != nil {
		return err
	}
	defer db.Close()
	store := model.NewProtocolInstanceStore(db)
	err = store.DeleteByProtocolAndChain(c.String("protocol"), c.String("chain"))
	if err != nil {
		return err
	}
	fmt.Printf("Deleted protocol %s on %s\n", c.String("protocol"), c.String("chain"))
	return nil
}

func protocolInstanceListHandler(c *cli.Context) error {
	db, err := model.DbConnect()
	if err != nil {
		return err
	}
	defer db.Close()
	store := model.NewProtocolInstanceStore(db)
	instances, err := store.All()
	if err != nil {
		return err
	}
	for _, instance := range instances {
		fmt.Printf("%s on %s @ %s\n", instance.Protocol.Name, instance.Chain.Name, instance.ContractAddress)
	}
	return nil
}

var registerProtocolInstanceCmd = &cli.Command{
	Name:   "register",
	Usage:  "registers a new protocol instance",
	Action: protocolInstanceRegisterHandler,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "chain",
			Usage:    "e.g polygon",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "protocol",
			Usage:    "e.g compound",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "contractAddress",
			Usage:    "e.g 0x1234",
			Required: true,
		},
	},
}

var deleteProtocolInstanceCmd = &cli.Command{
	Name:   "delete",
	Usage:  "deletes a protocol instance and its events",
	Action: protocolInstanceDeleteHandler,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "chain",
			Usage:    "e.g polygon",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "protocol",
			Usage:    "e.g compound",
			Required: true,
		},
	},
}

var listProtocolInstanceCmd = &cli.Command{
	Name:   "ls",
	Usage:  "lists all chain instances",
	Action: protocolInstanceListHandler,
}

var ProtocolInstanceCommand = &cli.Command{
	Name:  "protocolinstance",
	Usage: "registers and manages protocols and their ABIs",
	Flags: []cli.Flag{},
	Subcommands: []*cli.Command{

		registerProtocolInstanceCmd, deleteProtocolInstanceCmd, listProtocolInstanceCmd,
	},
}
