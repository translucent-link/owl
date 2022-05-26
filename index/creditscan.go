package index

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
	"github.com/translucent-link/owl/graph/model"
)

type ScanRequest struct {
	Client           *ethclient.Client
	Chain            *model.Chain
	Protocol         *model.Protocol
	ProtocolInstance *model.ProtocolInstance
	ScannableEvents  []*model.EventDefn
}

var ScanChannel chan ScanRequest

func init() {
	ScanChannel = make(chan ScanRequest)
	log.Println("Initialised  Scan channel")
	go func() {
		req := <-ScanChannel
		log.Println("Received Scan Request")
		ScanHistory(req.Client, req.Chain, req.Protocol, req.ProtocolInstance, req.ScannableEvents)
	}()
}

func grabUnpacker(protocol *model.Protocol) (Unpacker, error) {
	if protocol.Name == "Aave" {
		return UnpackAaveEvent, nil
	} else {
		return UnpackCompoundEvent, nil
	}
}

func ScanHistory(client *ethclient.Client, chain *model.Chain, protocol *model.Protocol, protocolInstance *model.ProtocolInstance, scannableEvents []*model.EventDefn) error {
	db, err := model.DbConnect()
	if err != nil {
		return errors.Wrap(err, "Unable to connect to DB and launch scan")
	}
	defer db.Close()
	stores := model.GenerateStores(db)

	unpacker, err := grabUnpacker(protocol)
	contractAbi, err := abi.JSON(strings.NewReader(protocol.Abi))
	if err != nil {
		return errors.Wrap(err, "Unable to parse ABI")
	}

	unknownTopics := []string{}
	currentBlock := big.NewInt(int64(protocolInstance.ScanStartBlock()))
	endBlock := big.NewInt(currentBlock.Int64())

	highestBlock, err := FindFirstBlock(client, 0)
	if err != nil {
		return errors.Wrap(err, "Finding first block")
	}
	toBlock := big.NewInt(highestBlock)
	step := big.NewInt(int64(chain.BlockFetchSize))

	for currentBlock.Int64() < toBlock.Int64() {

		endBlock.Add(endBlock, step)
		if endBlock.Int64() > toBlock.Int64() {
			endBlock = toBlock
		}

		log.Printf("Scanning from %d to %s", currentBlock, endBlock.String())

		topics := []common.Hash{}
		for _, event := range scannableEvents {
			topics = append(topics, common.HexToHash(event.TopicHashHex))
		}
		query := ethereum.FilterQuery{
			Addresses: []common.Address{common.HexToAddress(protocolInstance.ContractAddress)},
			FromBlock: currentBlock,
			ToBlock:   endBlock,
			Topics:    [][]common.Hash{topics},
		}

		logs, err := client.FilterLogs(context.Background(), query)
		if err != nil {
			log.Println(errors.Wrap(err, "Filtering logs"))
		} else {
			for _, vLog := range logs {
				found := false
				for _, eventDefn := range scannableEvents {
					if vLog.Topics[0].Hex() == eventDefn.TopicHashHex {
						found = true
						event, err := unpacker(contractAbi, protocolInstance, eventDefn, vLog)
						if err != nil {
							log.Println(errors.Wrap(err, "Unable to unpack event"))
						}

						// occuredAt =
						occuredAt, err := GetBlockTimestamp(client, big.NewInt(int64(vLog.BlockNumber)))
						if err != nil {
							log.Println(err)
						}

						if event.Borrowable != nil {
							borrower, err := stores.Account.FindOrCreateByAddress(event.Borrowable.GetBorrower().Hex())
							if err != nil {
								log.Println(errors.Wrap(err, "Unable find/create borrower for Borrow event"))
							}
							borrowToken, err := stores.Token.FindOrCreateByAddress(event.Borrowable.GetBorrowToken().Hex())
							if err != nil {
								log.Println(errors.Wrap(err, "Unable find/create token for Borrow event"))
							}
							_, err = stores.Event.StoreBorrowEvent(
								protocolInstance.ID,
								eventDefn.ID,
								vLog.TxHash.Hex(),
								int64(vLog.BlockNumber),
								occuredAt,
								borrower.ID,
								event.Borrowable.GetBorrowAmount(),
								borrowToken.ID)
							if err != nil {
								log.Println(errors.Wrap(err, "Unable store Borrow event"))
							}
						} else if event.Depositable != nil {
							depositor, err := stores.Account.FindOrCreateByAddress(event.Depositable.GetDepositor().Hex())
							if err != nil {
								log.Println(errors.Wrap(err, "Unable find/create borrower for Deposit event"))
							}
							depositToken, err := stores.Token.FindOrCreateByAddress(event.Depositable.GetDepositToken().Hex())
							if err != nil {
								log.Println(errors.Wrap(err, "Unable find/create token for Deposit event"))
							}
							_, err = stores.Event.StoreDepositEvent(
								protocolInstance.ID,
								eventDefn.ID,
								vLog.TxHash.Hex(),
								int64(vLog.BlockNumber),
								occuredAt,
								depositor.ID,
								event.Depositable.GetDepositAmount(),
								depositToken.ID)
							if err != nil {
								log.Println(errors.Wrap(err, "Unable store Deposit event"))
							}
						} else if event.Repayable != nil {
							borrower, err := stores.Account.FindOrCreateByAddress(event.Repayable.GetBorrower().Hex())
							if err != nil {
								log.Println(errors.Wrap(err, "Unable find/create borrower for Repay event"))
							}
							repayToken, err := stores.Token.FindOrCreateByAddress(event.Repayable.GetBorrowToken().Hex())
							if err != nil {
								log.Println(errors.Wrap(err, "Unable find/create token for Repay event"))
							}
							_, err = stores.Event.StoreRepayEvent(
								protocolInstance.ID,
								eventDefn.ID,
								vLog.TxHash.Hex(),
								int64(vLog.BlockNumber),
								occuredAt,
								borrower.ID,
								event.Repayable.GetRepayAmount(),
								repayToken.ID)
							if err != nil {
								log.Println(errors.Wrap(err, "Unable store Borrow event"))
							}
						} else if event.Liquidatable != nil {

							borrower, err := stores.Account.FindOrCreateByAddress(event.Liquidatable.GetBorrower().Hex())
							if err != nil {
								log.Println(errors.Wrap(err, "Unable find/create borrower for Liquidation event"))
							}
							liquidator, err := stores.Account.FindOrCreateByAddress(event.Liquidatable.GetLiquidator().Hex())
							if err != nil {
								log.Println(errors.Wrap(err, "Unable find/create borrower for Liquidation event"))
							}
							debtToken, err := stores.Token.FindOrCreateByAddress(event.Liquidatable.GetDebtToken().Hex())
							if err != nil {
								log.Println(errors.Wrap(err, "Unable find/create token for Repay event"))
							}
							collateralToken, err := stores.Token.FindOrCreateByAddress(event.Liquidatable.GetCollateralToken().Hex())
							if err != nil {
								log.Println(errors.Wrap(err, "Unable find/create token for Repay event"))
							}

							_, err = stores.Event.StoreLiquidationEvent(
								protocolInstance.ID,
								eventDefn.ID,
								vLog.TxHash.Hex(),
								int64(vLog.BlockNumber),
								occuredAt,
								borrower.ID,
								liquidator.ID,
								event.Liquidatable.GetRepayAmount(),
								event.Liquidatable.GetSeizeAmount(),
								debtToken.ID,
								collateralToken.ID)
						}

						// log.Println(vLog.Topics[0].Hex())
						// log.Println(event)
						break
					}
				}
				if !found {
					if indexOf(unknownTopics, vLog.Topics[0].Hex()) < 0 {
						log.Printf("UNK: %s", vLog.Topics[0].Hex())
						unknownTopics = append(unknownTopics, vLog.Topics[0].Hex())
					}

				}
			}
		}
		err = stores.ProtocolInstance.UpdateLastBlockRead(protocolInstance.ID, uint(endBlock.Int64()))
		if err != nil {
			return errors.Wrap(err, "Updating last block")
		}

		currentBlock = big.NewInt(endBlock.Int64())
	}
	log.Println("Uknown Topics...")
	for _, t := range unknownTopics {
		log.Println(t)
	}
	return nil
}

func ListenToEvents(client *ethclient.Client, address string) {
	contractAddress := common.HexToAddress(address)
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
	}

	logs := make(chan types.Log)
	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case vLog := <-logs:
			handleLogEvent(vLog)
		}
	}
}

func FindFirstBlock(client *ethclient.Client, days int) (int64, error) {

	currentTimestamp := time.Now().Unix()
	targetTimestamp := uint64(currentTimestamp - int64(days*86400))

	currentBlock, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return 0, err
	}

	currentBlockNumber := currentBlock.Number.Int64()
	if days > 0 {
		currentBlockTimestamp := currentBlock.Time
		currentStep := 1000

		for currentBlockTimestamp > targetTimestamp {
			currentBlock, err = client.HeaderByNumber(context.Background(), big.NewInt(currentBlockNumber-int64(currentStep)))
			if err != nil {
				return 0, err
			}
			currentBlockNumber = currentBlock.Number.Int64()

			fmt.Printf("%d %s\n", currentBlock.Number.Int64(), time.Unix(int64(currentBlockTimestamp), 0).String())
			if currentBlock.Time > (currentBlockTimestamp - 86400) {
				currentStep = currentStep * 10
			}
			currentBlockTimestamp = currentBlock.Time
		}
	}

	return currentBlockNumber, nil
}
