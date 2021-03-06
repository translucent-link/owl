package cmd

import (
	"log"

	"github.com/pkg/errors"
	"github.com/translucent-link/owl/graph/model"
	"github.com/translucent-link/owl/index"
	"github.com/urfave/cli/v2"
)

func scan(c *cli.Context) error {
	db, err := model.DbConnect()
	defer db.Close()
	stores := model.GenerateStores(db)

	chains, err := stores.Chain.All()
	if err != nil {
		return errors.Wrap(err, "Retrieving list of chains")
	}

	for _, chain := range chains {

		client, err := index.GetClient(chain.RPCURL)
		if err != nil {
			return errors.Wrap(err, "Retrieving EVM client")
		}

		protocols, err := stores.Protocol.AllByChain(chain.ID)
		if err != nil {
			return errors.Wrap(err, "Retrieving list of protocols")
		}
		for _, protocol := range protocols {
			protocolInstance, err := stores.ProtocolInstance.FindByProtocolIdAndChainId(protocol.ID, chain.ID)
			if err != nil {
				return errors.Wrap(err, "Retrieving list of protocol instances")
			}
			scannableEvents, err := stores.Protocol.AllEventsByProtocol(protocol.ID)
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
