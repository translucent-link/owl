package index

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Unpacker func(abi abi.ABI, le LoanEvent, data []byte) (Typable, error)

type Protocol struct {
	Name            string
	ContractAddress common.Address
	ScannableEvents []LoanEvent
	TopicHashes     []common.Hash
	ContractABI     abi.ABI
	Unpacker        Unpacker
}

type Typable interface {
	Type() string
}

type LoanEvent struct {
	TopicName    string
	TopicHash    common.Hash
	TopicHashHex string
}

func (le LoanEvent) String() string {
	return fmt.Sprintf("%s @ %s", le.TopicName, le.TopicHashHex)
}

func debugPrint(msg string, value interface{}) {
	fmt.Printf("%s %#v\n", msg, value)
}

func handleLogEvent(vLog types.Log) {
	debugPrint("Log Event", vLog)
}

func GetClient(ethURL string) (*ethclient.Client, error) {
	return ethclient.Dial(ethURL)
}
