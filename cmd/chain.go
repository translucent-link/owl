package cmd

import (
	"fmt"

	model "github.com/translucent-link/owl/graph/model"

	"github.com/urfave/cli/v2"
)

func chainRegisterHandler(c *cli.Context) error {
	db, err := model.DbConnect()
	if err != nil {
		return err
	}
	defer db.Close()
	store := model.NewChainStore(db)
	newChain, err := store.CreateChain(model.NewChain{
		Name:           c.String("name"),
		RPCURL:         c.String("rpcURL"),
		BlockFetchSize: c.Int("blockFetchSize"),
	})
	if err != nil {
		return err
	}
	fmt.Printf("Registered new chain: %s\n", newChain.Name)
	return nil
}

func chainDeleteHandler(c *cli.Context) error {
	db, err := model.DbConnect()
	if err != nil {
		return err
	}
	defer db.Close()
	store := model.NewChainStore(db)
	err = store.DeleteByName(c.String("name"))
	if err != nil {
		return err
	}
	fmt.Printf("Deleted chain: %s\n", c.String("name"))
	return nil
}

func chainListHandler(c *cli.Context) error {
	db, err := model.DbConnect()
	if err != nil {
		return err
	}
	defer db.Close()
	store := model.NewChainStore(db)
	chains, err := store.All()
	if err != nil {
		return err
	}
	for _, chain := range chains {
		fmt.Printf("%s\n", chain.Name)
	}
	return nil
}

var registerChainCmd = &cli.Command{
	Name:   "register",
	Usage:  "registers a new chain",
	Action: chainRegisterHandler,
	Flags: []cli.Flag{
		&cli.IntFlag{
			Name:  "blockFetchSize",
			Usage: "e.g. 10000",
		},
		&cli.StringFlag{
			Name:     "name",
			Usage:    "e.g polygon",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "rpcURL",
			Usage:    "e.g wss://polygon-mainnet.g.alchemy.com/v2/<APIKEY>",
			Required: true,
		},
	},
}

var deleteChainCmd = &cli.Command{
	Name:   "delete",
	Usage:  "deletes a chain and all associated protocol instances and events",
	Action: chainDeleteHandler,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "name",
			Usage:    "e.g polygon",
			Required: true,
		},
	},
}

var listChainCmd = &cli.Command{
	Name:   "ls",
	Usage:  "lists all chains",
	Action: chainListHandler,
}

var ChainCommand = &cli.Command{
	Name:  "chain",
	Usage: "registers and manages chains",
	Flags: []cli.Flag{},
	Subcommands: []*cli.Command{

		registerChainCmd, deleteChainCmd, listChainCmd,
	},
}
