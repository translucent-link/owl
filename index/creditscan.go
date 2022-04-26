package index

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func ScanHistory(client *ethclient.Client, abiPath string, fromBlock *big.Int, toBlock *big.Int) {
	protocol, err := GetAEthDefn(abiPath)
	if err != nil {
		log.Fatal(err)
	}

	currentBlock := fromBlock
	endBlock := big.NewInt(0)
	step := big.NewInt(1000)
	for currentBlock.Int64() < toBlock.Int64() {

		endBlock.Add(currentBlock, step)
		if endBlock.Int64() > toBlock.Int64() {
			endBlock = toBlock
		}

		log.Printf("Scanning from %s to %s", currentBlock.String(), endBlock.String())
		query := ethereum.FilterQuery{
			Addresses: []common.Address{protocol.ContractAddress},
			FromBlock: currentBlock,
			ToBlock:   endBlock,
			Topics:    [][]common.Hash{protocol.TopicHashes},
		}

		logs, err := client.FilterLogs(context.Background(), query)
		if err != nil {
			log.Fatal(err)
		}

		for _, vLog := range logs {
			debugPrint("\n", vLog)
			for _, loanEvent := range protocol.ScannableEvents {
				if vLog.Topics[0].Hex() == loanEvent.TopicHashHex {

					fmt.Println(loanEvent.String())
					event, err := protocol.Unpacker(protocol.ContractABI, loanEvent, vLog.Data)
					if err != nil {
						log.Fatal(err)
					}

					log.Println(event)
					break
				}
			}
		}

		currentBlock = endBlock
	}

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
