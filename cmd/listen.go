package cmd

import (
	"github.com/translucent-link/owl/index"
	"github.com/urfave/cli/v2"
)

func listen(c *cli.Context) error {

	ethURL := c.String("ethURL")
	contractAddress := c.String("contractAddress")

	client, err := index.GetClient(ethURL)
	if err != nil {
		return err
	}

	index.ListenToEvents(client, contractAddress)
	return nil
}

var ListenCommand = &cli.Command{
	Name:   "listen",
	Usage:  "listens to events on the blockchain",
	Action: listen,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "ethURL",
			Aliases:  []string{"u"},
			Usage:    "wss:// or https:// URL pointing to blockchain node",
			Required: true,
			EnvVars:  []string{"ETH_URL"},
		},
		&cli.StringFlag{
			Name:     "contractAddress",
			Aliases:  []string{"c"},
			Usage:    "hexedecimal string address pointing to contract",
			Required: true,
			EnvVars:  []string{"LISTEN_CONTRACT_ADDR"},
		},
	},
}
