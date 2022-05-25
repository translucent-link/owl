package cmd

import (
	"log"

	"github.com/pkg/errors"
	"github.com/translucent-link/owl/graph/model"
	"github.com/translucent-link/owl/index"
	"github.com/urfave/cli/v2"
)

func scan(c *cli.Context) error {

	chainStore, protocolStore, protocolInstanceStore, err := model.Stores()

	chains, err := chainStore.All()
	if err != nil {
		return errors.Wrap(err, "Retrieving list of chains")
	}

	for _, chain := range chains {

		client, err := index.GetClient(chain.RPCURL)
		if err != nil {
			return errors.Wrap(err, "Retrieving EVM client")
		}

		protocols, err := protocolStore.AllByChain(chain.ID)
		if err != nil {
			return errors.Wrap(err, "Retrieving list of protocols")
		}
		for _, protocol := range protocols {
			protocolInstance, err := protocolInstanceStore.FindByProtocolIdAndChainId(protocol.ID, chain.ID)
			if err != nil {
				return errors.Wrap(err, "Retrieving list of protocol instances")
			}
			scannableEvents, err := protocolStore.AllEventsByProtocol(protocol.ID)
			if err != nil {
				return errors.Wrap(err, "Retrieving list of scannable events")
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
