package cmd

import (
	"fmt"

	"github.com/translucent-link/owl/index"
	"github.com/urfave/cli/v2"
)

func blk(c *cli.Context) error {
	days := c.Int("days")
	ethURL := c.String("ethURL")
	client, err := index.GetClient(ethURL)
	if err != nil {
		return err
	}
	blk, err := index.FindFirstBlock(client, days)
	if err != nil {
		return err
	}

	fmt.Println(blk)
	return nil
}

var BlkCommand = &cli.Command{
	Name:   "blk",
	Usage:  "helps discover particular blocks",
	Action: blk,
	Flags: []cli.Flag{
		&cli.IntFlag{
			Name:  "days",
			Usage: "how many days to travel back",
		},
		&cli.StringFlag{
			Name:     "ethURL",
			Aliases:  []string{"u"},
			Usage:    "wss:// or https:// URL pointing to blockchain node",
			Required: true,
			EnvVars:  []string{"ETH_URL"},
		},
	},
}
