package cmd

import (
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/urfave/cli/v2"
)

func topicHashCmdHandler(c *cli.Context) error {
	topicDefn := c.Args().Get(0)
	topicSignature := []byte(topicDefn)
	topicHash := crypto.Keccak256Hash(topicSignature)

	fmt.Println(topicHash)
	return nil
}

var topicHashCmd = &cli.Command{
	Name:   "topicHash",
	Usage:  "helps with ABI-encoding",
	Action: topicHashCmdHandler,
}

var AbiCommand = &cli.Command{
	Name:  "abi",
	Usage: "converts an ABI-event definition into a topic hash",
	Flags: []cli.Flag{},
	Subcommands: []*cli.Command{

		topicHashCmd,
	},
}
