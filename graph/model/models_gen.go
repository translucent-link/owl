// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"io"
	"strconv"
	"time"
)

type AnyEvent interface {
	IsAnyEvent()
}

type Event interface {
	IsEvent()
}

type Account struct {
	ID      int        `json:"id"`
	Address string     `json:"address"`
	Events  []AnyEvent `json:"events"`
}

type BorrowEvent struct {
	ID             int       `json:"id"`
	Type           EventType `json:"type"`
	Txhash         string    `json:"txhash"`
	Blocknumber    int       `json:"blocknumber"`
	Index          int       `json:"index"`
	OccuredAt      time.Time `json:"occuredAt"`
	Borrower       *Account  `json:"borrower"`
	AmountBorrowed int       `json:"amountBorrowed"`
	Token          *Token    `json:"token"`
}

func (BorrowEvent) IsEvent()    {}
func (BorrowEvent) IsAnyEvent() {}

type Chain struct {
	ID             int         `json:"id"`
	Name           string      `json:"name"`
	RPCURL         string      `json:"rpcUrl"`
	BlockFetchSize int         `json:"blockFetchSize"`
	Protocols      []*Protocol `json:"protocols"`
	Tokens         []*Token    `json:"tokens"`
}

type DepositEvent struct {
	ID              int       `json:"id"`
	Type            EventType `json:"type"`
	Txhash          string    `json:"txhash"`
	Blocknumber     int       `json:"blocknumber"`
	Index           int       `json:"index"`
	OccuredAt       time.Time `json:"occuredAt"`
	Depositor       *Account  `json:"depositor"`
	AmountDeposited int       `json:"amountDeposited"`
	Token           *Token    `json:"token"`
}

func (DepositEvent) IsEvent()    {}
func (DepositEvent) IsAnyEvent() {}

type EventDefn struct {
	ID           int    `json:"id"`
	TopicName    string `json:"topicName"`
	TopicHashHex string `json:"topicHashHex"`
	AbiSignature string `json:"abiSignature"`
}

type LiquidationEvent struct {
	ID              int       `json:"id"`
	Type            EventType `json:"type"`
	Txhash          string    `json:"txhash"`
	Blocknumber     int       `json:"blocknumber"`
	Index           int       `json:"index"`
	OccuredAt       time.Time `json:"occuredAt"`
	Borrower        *Account  `json:"borrower"`
	Liquidator      *Account  `json:"liquidator"`
	AmountRepayed   int       `json:"amountRepayed"`
	AmountSeized    int       `json:"amountSeized"`
	CollateralToken *Token    `json:"collateralToken"`
	DebtToken       *Token    `json:"debtToken"`
}

func (LiquidationEvent) IsEvent()    {}
func (LiquidationEvent) IsAnyEvent() {}

type NewChain struct {
	Name           string `json:"name"`
	RPCURL         string `json:"rpcUrl"`
	BlockFetchSize int    `json:"blockFetchSize"`
}

type NewEventDefn struct {
	Protocol     string `json:"protocol"`
	TopicName    string `json:"topicName"`
	AbiSignature string `json:"abiSignature"`
}

type NewProtocol struct {
	Name string `json:"name"`
	Abi  string `json:"abi"`
}

type NewProtocolInstance struct {
	Protocol         string `json:"protocol"`
	Chain            string `json:"chain"`
	ContractAddress  string `json:"contractAddress"`
	FirstBlockToRead int    `json:"firstBlockToRead"`
}

type NewScan struct {
	Protocol string `json:"protocol"`
	Chain    string `json:"chain"`
}

type Protocol struct {
	ID              int          `json:"id"`
	Name            string       `json:"name"`
	Abi             string       `json:"abi"`
	ScannableEvents []*EventDefn `json:"scannableEvents"`
}

type ProtocolInstance struct {
	ID               int       `json:"id"`
	Protocol         *Protocol `json:"protocol"`
	Chain            *Chain    `json:"chain"`
	ContractAddress  string    `json:"contractAddress"`
	FirstBlockToRead int       `json:"firstBlockToRead"`
	LastBlockRead    int       `json:"lastBlockRead"`
}

type RepayEvent struct {
	ID            int       `json:"id"`
	Type          EventType `json:"type"`
	Txhash        string    `json:"txhash"`
	Blocknumber   int       `json:"blocknumber"`
	Index         int       `json:"index"`
	OccuredAt     time.Time `json:"occuredAt"`
	Borrower      *Account  `json:"borrower"`
	AmountRepayed int       `json:"amountRepayed"`
	Token         *Token    `json:"token"`
}

func (RepayEvent) IsEvent()    {}
func (RepayEvent) IsAnyEvent() {}

type Token struct {
	ID       int     `json:"id"`
	Address  string  `json:"address"`
	Name     *string `json:"name"`
	Ticker   *string `json:"ticker"`
	Decimals int     `json:"decimals"`
}

type TokenInfo struct {
	Address  string `json:"address"`
	Name     string `json:"name"`
	Ticker   string `json:"ticker"`
	Chain    string `json:"chain"`
	Decimals int    `json:"decimals"`
}

type EventType string

const (
	EventTypeBorrow      EventType = "Borrow"
	EventTypeRepay       EventType = "Repay"
	EventTypeLiquidation EventType = "Liquidation"
	EventTypeDeposit     EventType = "Deposit"
)

var AllEventType = []EventType{
	EventTypeBorrow,
	EventTypeRepay,
	EventTypeLiquidation,
	EventTypeDeposit,
}

func (e EventType) IsValid() bool {
	switch e {
	case EventTypeBorrow, EventTypeRepay, EventTypeLiquidation, EventTypeDeposit:
		return true
	}
	return false
}

func (e EventType) String() string {
	return string(e)
}

func (e *EventType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = EventType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid EventType", str)
	}
	return nil
}

func (e EventType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
