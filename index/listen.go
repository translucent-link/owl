package index

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/translucent-link/owl/graph/model"
	"github.com/translucent-link/owl/utils"

	_ "github.com/lib/pq"
	_ "github.com/mattes/migrate/database/postgres"
	_ "github.com/mattes/migrate/source/file"
)

func init() {
	_ = godotenv.Load()

	utils.SetupDatabase()

	stores, err := model.NewStores()
	if err != nil {
		log.Fatal(errors.Wrap(err, "Unable to connect to DB whilst setting up listeners"))
	}
	defer stores.Close()

	chains, err := stores.Chain.All()
	if err != nil {
		log.Fatal(err, "Retrieving list of chains")
	}

	for _, chain := range chains {

		if strings.HasPrefix(chain.RPCURL, "ws") {
			protocols, err := stores.Protocol.AllByChain(chain.ID)
			if err != nil {
				log.Fatal(errors.Wrap(err, "Retrieving list of protocols"))
			}
			for _, protocol := range protocols {
				protocolInstance, err := stores.ProtocolInstance.FindByProtocolIdAndChainId(protocol.ID, chain.ID)
				if err != nil {
					log.Fatal(errors.Wrap(err, "Retrieving list of protocol instances"))
				}
				scannableEvents, err := stores.Protocol.AllEventsByProtocol(protocol.ID)
				if err != nil {
					log.Fatal(errors.Wrap(err, "Retrieving list of scannable events"))
				}
				log.Printf("Listening to %s on %s", protocol.Name, chain.Name)
				go ListenToEvents(chain, protocol, protocolInstance, scannableEvents)
			}
		} else {
			log.Printf("Skipping %s. Requires websockets", chain.Name)
		}
	}

}

func ListenToEvents(chain *model.Chain, protocol *model.Protocol, protocolInstance *model.ProtocolInstance, scannableEvents []*model.EventDefn) {
	topics := []common.Hash{}
	for _, event := range scannableEvents {
		topics = append(topics, common.HexToHash(event.TopicHashHex))
	}
	contractAddress := common.HexToAddress(protocolInstance.ContractAddress)
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
		Topics:    [][]common.Hash{topics},
	}

	var client *ethclient.Client
	var err error
	for {
		client, err = chain.EthClient()
		if err != nil {
			log.Println(errors.Wrapf(err, "Unable to connect to %s", chain.RPCURL).Error())
			break
		}
		logs := make(chan types.Log)
		sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
		if err != nil {
			log.Fatal(errors.Wrap(err, "Unable to subscribe to topics"))
		}

		for {
			select {
			case err := <-sub.Err():
				log.Println(errors.Wrap(err, "Received err whilst listening to events"))
			case vLog := <-logs:
				handleListenLogEvent(vLog, chain, protocol, protocolInstance, scannableEvents)
			}
		}

	}
}

func handleListenLogEvent(vLog types.Log, chain *model.Chain, protocol *model.Protocol, protocolInstance *model.ProtocolInstance, scannableEvents []*model.EventDefn) {
	db, err := model.DbConnect()
	if err != nil {
		log.Fatal(errors.Wrap(err, "Unable to connect to DB whilst handling log event"))
	}
	defer db.Close()
	stores := model.GenerateStores(db)

	contractAbi, err := abi.JSON(strings.NewReader(protocol.Abi))
	if err != nil {
		log.Println(errors.Wrap(err, "Unable to parse ABI"))
	}
	unpacker, err := grabUnpacker(protocol)

	for _, eventDefn := range scannableEvents {
		if vLog.Topics[0].Hex() == eventDefn.TopicHashHex {
			event, err := unpacker(contractAbi, protocolInstance, eventDefn, vLog)

			if err != nil {
				log.Println(errors.Wrap(err, "Unable to unpack event"))
			}

			occuredAt := time.Now()
			log.Printf("%s: %s: %s", chain.Name, protocol.Name, vLog.TxHash.Hex())
			if event.Borrowable != nil {
				log.Println(event.Borrowable)
				borrower, err := stores.Account.FindOrCreateByAddress(event.Borrowable.GetBorrower().Hex())
				if err != nil {
					log.Println(errors.Wrap(err, "Unable find/create borrower for Borrow event"))
				}
				borrowToken, err := stores.Token.FindOrCreateByAddress(event.Borrowable.GetBorrowToken().Hex(), chain.ID)
				if err != nil {
					log.Println(errors.Wrap(err, "Unable find/create token for Borrow event"))
				}
				_, err = stores.Event.StoreBorrowEvent(
					protocolInstance.ID,
					eventDefn.ID,
					vLog.TxHash.Hex(),
					int64(vLog.BlockNumber),
					int(vLog.Index),
					occuredAt,
					borrower.ID,
					event.Borrowable.GetBorrowAmount(),
					borrowToken.ID)
				if err != nil && !isDuplicateError(err) {
					log.Println(errors.Wrapf(err, "Unable store Borrow event %s on PI:%d identified by %s", vLog.TxHash.Hex(), protocolInstance.ID, eventDefn.TopicHashHex))
				}
			} else if event.Depositable != nil {
				log.Println(event.Depositable)
				depositor, err := stores.Account.FindOrCreateByAddress(event.Depositable.GetDepositor().Hex())
				if err != nil {
					log.Println(errors.Wrap(err, "Unable find/create borrower for Deposit event"))
				}
				depositToken, err := stores.Token.FindOrCreateByAddress(event.Depositable.GetDepositToken().Hex(), chain.ID)
				if err != nil {
					log.Println(errors.Wrap(err, "Unable find/create token for Deposit event"))
				}
				_, err = stores.Event.StoreDepositEvent(
					protocolInstance.ID,
					eventDefn.ID,
					vLog.TxHash.Hex(),
					int64(vLog.BlockNumber),
					int(vLog.Index),
					occuredAt,
					depositor.ID,
					event.Depositable.GetDepositAmount(),
					depositToken.ID)
				if err != nil && !isDuplicateError(err) {
					log.Println(errors.Wrapf(err, "Unable store Deposit event %s on PI:%d identified by %s", vLog.TxHash.Hex(), protocolInstance.ID, eventDefn.TopicHashHex))
				}
			} else if event.Repayable != nil {
				log.Println(event.Repayable)
				borrower, err := stores.Account.FindOrCreateByAddress(event.Repayable.GetBorrower().Hex())
				if err != nil {
					log.Println(errors.Wrap(err, "Unable find/create borrower for Repay event"))
				}
				repayToken, err := stores.Token.FindOrCreateByAddress(event.Repayable.GetBorrowToken().Hex(), chain.ID)
				if err != nil {
					log.Println(errors.Wrap(err, "Unable find/create token for Repay event"))
				}
				_, err = stores.Event.StoreRepayEvent(
					protocolInstance.ID,
					eventDefn.ID,
					vLog.TxHash.Hex(),
					int64(vLog.BlockNumber),
					int(vLog.Index),
					occuredAt,
					borrower.ID,
					event.Repayable.GetRepayAmount(),
					repayToken.ID)
				if err != nil && !isDuplicateError(err) {
					log.Println(errors.Wrapf(err, "Unable store Repay event %s on PI:%d identified by %s", vLog.TxHash.Hex(), protocolInstance.ID, eventDefn.TopicHashHex))
				}
			} else if event.Liquidatable != nil {
				log.Println(event.Liquidatable)

				borrower, err := stores.Account.FindOrCreateByAddress(event.Liquidatable.GetBorrower().Hex())
				if err != nil {
					log.Println(errors.Wrap(err, "Unable find/create borrower for Liquidation event"))
				}
				liquidator, err := stores.Account.FindOrCreateByAddress(event.Liquidatable.GetLiquidator().Hex())
				if err != nil {
					log.Println(errors.Wrap(err, "Unable find/create borrower for Liquidation event"))
				}
				debtToken, err := stores.Token.FindOrCreateByAddress(event.Liquidatable.GetDebtToken().Hex(), chain.ID)
				if err != nil {
					log.Println(errors.Wrap(err, "Unable find/create token for Repay event"))
				}
				collateralToken, err := stores.Token.FindOrCreateByAddress(event.Liquidatable.GetCollateralToken().Hex(), chain.ID)
				if err != nil {
					log.Println(errors.Wrap(err, "Unable find/create token for Repay event"))
				}

				_, err = stores.Event.StoreLiquidationEvent(
					protocolInstance.ID,
					eventDefn.ID,
					vLog.TxHash.Hex(),
					int64(vLog.BlockNumber),
					int(vLog.Index),
					occuredAt,
					borrower.ID,
					liquidator.ID,
					event.Liquidatable.GetRepayAmount(),
					event.Liquidatable.GetSeizeAmount(),
					debtToken.ID,
					collateralToken.ID)
				if err != nil && !isDuplicateError(err) {
					log.Println(errors.Wrapf(err, "Unable store Liquidation event %s on PI:%d identified by %s", vLog.TxHash.Hex(), protocolInstance.ID, eventDefn.TopicHashHex))
				}

			}

			break
		}
	}
}
