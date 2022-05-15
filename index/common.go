package index

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/translucent-link/owl/graph/model"
)

type Unpacker func(abi abi.ABI, eventDefn *model.EventDefn, log types.Log) (Typable, error)

type Chain struct {
	RpcURL         string
	BlockFetchSize int
}

// type Protocol struct {
// 	Name            string
// 	ContractAddress common.Address
// 	ScannableEvents []LoanEvent
// 	TopicHashes     []common.Hash
// 	ContractABI     abi.ABI
// 	Unpacker        Unpacker
// }

// type ProtocolInstance struct {
// 	Protocol
// 	Chain
// }

type Typable interface {
	Type() string
}

type Transactional interface {
	Txn() string
}

// type LoanEvent struct {
// 	TopicName    string
// 	TopicHash    common.Hash
// 	TopicHashHex string
// }

// func (le LoanEvent) String() string {
// 	return fmt.Sprintf("%s @ %s", le.TopicName, le.TopicHashHex)
// }

func debugPrint(msg string, value interface{}) {
	fmt.Printf("%s %#v\n", msg, value)
}

func handleLogEvent(vLog types.Log) {
	debugPrint("Log Event", vLog)
}

func GetClient(ethURL string) (*ethclient.Client, error) {
	return ethclient.Dial(ethURL)
}

func indexOf(values []string, tgt string) int {
	for i, v := range values {
		if v == tgt {
			return i
		}
	}
	return -1
}
