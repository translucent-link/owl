package cmd

import (
	"math/big"
	"strconv"

	"github.com/translucent-link/owl/index"
	"github.com/urfave/cli/v2"
)

func scan(c *cli.Context) error {
	fromBlockStr := c.Args().Get(0)
	toBlockStr := c.Args().Get(1)

	ethURL := c.String("ethURL")
	abiPath := c.String("abiPath")
	client, err := index.GetClient(ethURL)
	if err != nil {
		return err
	}

	fromBlock, err := strconv.Atoi(fromBlockStr)
	if err != nil {
		return err
	}
	toBlock, err := strconv.Atoi(toBlockStr)
	if err != nil {
		return err
	}

	index.ScanHistory(client, abiPath, big.NewInt(int64(fromBlock)), big.NewInt(int64(toBlock)))

	return nil
}

var ScanCommand = &cli.Command{
	Name:   "scan",
	Usage:  "helps discover particular blocks",
	Action: scan,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "ethURL",
			Aliases:  []string{"u"},
			Usage:    "wss:// or https:// URL pointing to blockchain node",
			Required: true,
			EnvVars:  []string{"ETH_URL"},
		},
		&cli.StringFlag{
			Name:     "abiPath",
			Aliases:  []string{"a"},
			Usage:    "subfolder containing abi files",
			Required: true,
			EnvVars:  []string{"ABI_PATH"},
		},
	},
}
