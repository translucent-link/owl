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
	"github.com/translucent-link/owl/graph/model"
)

func grabUnpacker(protocol *model.Protocol) (Unpacker, error) {
	if protocol.Name == "Aave" {
		return UnpackAaveEvent, nil
	} else {
		return UnpackCompoundEvent, nil
	}
}

func ScanHistory(client *ethclient.Client, chain *model.Chain, protocol *model.Protocol, protocolInstance *model.ProtocolInstance, scannableEvents []*model.EventDefn) error {

	protocolInstanceStore, err := model.NewProtocolInstanceStore()
	if err != nil {
		return err
	}
	unpacker, err := grabUnpacker(protocol)
	contractAbi, err := abi.JSON(strings.NewReader(protocol.Abi))
	if err != nil {
		return err
	}

	unknownTopics := []string{}
	currentBlock := big.NewInt(int64(protocolInstance.LastBlockRead))
	endBlock := big.NewInt(currentBlock.Int64())

	highestBlock, err := FindFirstBlock(client, 0)
	if err != nil {
		return err
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
			log.Println(err)
		} else {
			for _, vLog := range logs {
				found := false
				for _, eventDefn := range scannableEvents {
					if vLog.Topics[0].Hex() == eventDefn.TopicHashHex {
						found = true
						event, err := unpacker(contractAbi, eventDefn, vLog)
						if err != nil {
							log.Fatal(err)
						}

						log.Println(vLog.Topics[0].Hex())
						log.Println(event)
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
		protocolInstanceStore.UpdateLastBlockRead(protocolInstance.ID, uint(endBlock.Int64()))

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
