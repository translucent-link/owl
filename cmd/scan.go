package cmd

import (
	"log"

	"github.com/translucent-link/owl/graph/model"
	"github.com/translucent-link/owl/index"
	"github.com/urfave/cli/v2"
)

func scan(c *cli.Context) error {

	chainStore, protocolStore, protocolInstanceStore, err := model.Stores()

	chains, err := chainStore.All()
	if err != nil {
		return err
	}

	for _, chain := range chains {

		client, err := index.GetClient(chain.RPCURL)
		if err != nil {
			return err
		}

		protocols, err := protocolStore.AllByChain(chain.ID)
		if err != nil {
			return err
		}
		for _, protocol := range protocols {
			protocolInstance, err := protocolInstanceStore.FindByProtocolIdAndChainId(chain.ID, protocol.ID)
			if err != nil {
				return err
			}
			scannableEvents, err := protocolStore.AllEventsByProtocol(protocol.ID)
			if err != nil {
				return err
			}
			log.Printf("Scanning %s on %s", protocol.Name, chain.Name)
			index.ScanHistory(client, chain, protocol, protocolInstance, scannableEvents)
		}
	}

	return nil
}

var ScanCommand = &cli.Command{
	Name:   "scan",
	Usage:  "helps discover particular blocks",
	Action: scan,
}
