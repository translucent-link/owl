package cmd

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/translucent-link/owl/index"
	"github.com/urfave/cli/v2"
)

func lookupEnvVarBasedOnChain(chain string) string {
	switch chain {
	case "ethereum":
		return "ETH_URL"
	case "polygon":
		return "POLYGON_URL"
	case "sepolia":
		return "SEPOLIA_URL"
	case "arbitrum":
		return "ARBITRUM_URL"
	case "arbitrum_sepolia":
		return "ARBITRUM_SEPOLIA_URL"
	case "avalanche":
		return "AVALANCHE_URL"
	default:
		return ""
	}
}

func blk(c *cli.Context) error {
	days := c.Int("days")
	chain := c.String("chain")
	rpcURLVar := lookupEnvVarBasedOnChain(chain)
	client, err := index.GetClient(os.Getenv(rpcURLVar))
	if err != nil {
		return errors.Wrapf(err, "Please ensure you have set the %s environment variable", rpcURLVar)
	}
	blk, err := index.FindFirstBlock(client, days)
	if err != nil {
		return errors.Wrapf(err, "Please ensure you have set the %s environment variable", rpcURLVar)
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
			Name:     "chain",
			Aliases:  []string{"c"},
			Usage:    "ethereum, polygon, avalanche, sepolia, arbitrum, or arbitrum_sepolia",
			Required: true,
		},
	},
}
