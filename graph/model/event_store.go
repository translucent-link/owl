package model

import (
	"log"
	"math/big"
	"time"

	"github.com/jmoiron/sqlx"
)

type EventStore struct {
	db *sqlx.DB
}

func NewEventStore(db *sqlx.DB) *EventStore {
	return &EventStore{db: db}
}

func (s *EventStore) AllByAccount(accountId int) ([]AllEvent, error) {
	events := []AllEvent{}
	sql := `SELECT e.*
	FROM events e
	WHERE e.borrowerAccountId = $1 OR e.repayerAccountId = $1 OR e.depositorAccountId = $1 OR e.liquidatorAccountId = $1
	ORDER BY blocknumber ASC;`
	err := s.db.Select(&events, sql, accountId)
	return events, err
}

func (s *EventStore) StoreBorrowEvent(protocolInstanceId int, eventDefinitionId int, txHash string, blockNumber int64, occuredAt time.Time, borrowerId int, borrowAmount *big.Int, borrowTokenId int) (int, error) {
	var insertedId int
	err := s.db.QueryRowx(
		"insert into events (type, protocolInstanceId, eventDefinitionId, txHash, blockNumber, occuredAt, borrowerAccountId, amountBorrowed, borrowTokenId) values ($1,$2,$3,$4,$5,$6,$7,$8,$9) returning id",
		"Borrow", protocolInstanceId, eventDefinitionId, txHash, blockNumber, occuredAt, borrowerId, borrowAmount.String(), borrowTokenId).Scan(&insertedId)
	return insertedId, err
}

func (s *EventStore) StoreRepayEvent(protocolInstanceId int, eventDefinitionId int, txHash string, blockNumber int64, occuredAt time.Time, repayerAccountId int, repayAmount *big.Int, repayTokenId int) (int, error) {
	var insertedId int
	err := s.db.QueryRowx(
		"insert into events (type, protocolInstanceId, eventDefinitionId, txHash, blockNumber, occuredAt, repayerAccountId, amountRepayed, repayTokenId) values ($1,$2,$3,$4,$5,$6,$7,$8,$9) returning id",
		"Repay", protocolInstanceId, eventDefinitionId, txHash, blockNumber, occuredAt, repayerAccountId, repayAmount.String(), repayTokenId).Scan(&insertedId)
	return insertedId, err
}

func (s *EventStore) StoreLiquidationEvent(protocolInstanceId int, eventDefinitionId int, txHash string, blockNumber int64, occuredAt time.Time, borrowerAccountId int, liquidatorAccountId int, amountRepayed *big.Int, amountSeized *big.Int, debtTokenId int, collateralTokenId int) (int, error) {
	var insertedId int
	err := s.db.QueryRowx(
		"insert into events (type, protocolInstanceId, eventDefinitionId, txHash, blockNumber, occuredAt, borrowerAccountId, liquidatorAccountId, amountRepayed, amountSeized, debtTokenId, collateralTokenId) values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12) returning id",
		"Liquidation", protocolInstanceId, eventDefinitionId, txHash, blockNumber, occuredAt, borrowerAccountId, liquidatorAccountId, amountRepayed.String(), amountSeized.String(), debtTokenId, collateralTokenId).Scan(&insertedId)
	return insertedId, err
}

func (s *EventStore) StoreDepositEvent(protocolInstanceId int, eventDefinitionId int, txHash string, blockNumber int64, occuredAt time.Time, depositorAccountId int, amountDeposited *big.Int, depositTokenId int) (int, error) {
	var insertedId int
	err := s.db.QueryRowx(
		"insert into events (type, protocolInstanceId, eventDefinitionId, txHash, blockNumber, occuredAt, depositorAccountId, amountDeposited, depositTokenId) values ($1,$2,$3,$4,$5,$6,$7,$8,$9) returning id",
		"Deposit", protocolInstanceId, eventDefinitionId, txHash, blockNumber, occuredAt, depositorAccountId, amountDeposited.String(), depositTokenId).Scan(&insertedId)
	return insertedId, err
}

type AllEvent struct {
	// Common
	ID                 int
	Type               EventType
	ProtocolInstanceId int
	EventDefinitionId  int
	Txhash             string
	Blocknumber        int
	OccuredAt          time.Time

	// Borrow
	DepositorAccountId *int
	AmountDeposited    *int
	DepositTokenId     *int

	// Borrow
	BorrowerAccountId *int
	AmountBorrowed    *int
	BorrowTokenId     *int

	// Repay
	RepayerAccountId *int
	AmountRepayed    *int
	RepayTokenId     *int

	// Liquidation
	LiquidatorAccountId *int
	AmountSeized        *int
	CollateralTokenId   *int
	DebtTokenId         *int
}

func (e AllEvent) AnyEvent(accountStore *AccountStore, tokenStore *TokenStore) (AnyEvent, error) {
	if accountStore == nil {
		log.Panic("AccountStore is nil")
	}
	if tokenStore == nil {
		log.Panic("TokenStore is nil")
	}
	if e.Type == EventTypeDeposit {
		depositor, err := accountStore.FindById(*e.DepositorAccountId)
		if err != nil {
			return DepositEvent{}, err
		}
		token, err := tokenStore.FindById(*e.DepositTokenId)
		if err != nil {
			return DepositEvent{}, err
		}
		return DepositEvent{e.ID, e.Type, e.Txhash, e.Blocknumber, e.OccuredAt, depositor, *e.AmountDeposited, token}, nil
	} else if e.Type == EventTypeBorrow {
		borrower, err := accountStore.FindById(*e.BorrowerAccountId)
		if err != nil {
			return BorrowEvent{}, err
		}
		token, err := tokenStore.FindById(*e.BorrowTokenId)
		if err != nil {
			return BorrowEvent{}, err
		}
		return BorrowEvent{e.ID, e.Type, e.Txhash, e.Blocknumber, e.OccuredAt, borrower, *e.AmountBorrowed, token}, nil
	} else if e.Type == EventTypeRepay {
		repayer, err := accountStore.FindById(*e.RepayerAccountId)
		if err != nil {
			return RepayEvent{}, err
		}
		token, err := tokenStore.FindById(*e.RepayTokenId)
		if err != nil {
			return RepayEvent{}, err
		}

		return RepayEvent{e.ID, e.Type, e.Txhash, e.Blocknumber, e.OccuredAt, repayer, *e.AmountRepayed, token}, nil
	} else if e.Type == EventTypeLiquidation {
		borrower, err := accountStore.FindById(*e.BorrowerAccountId)
		if err != nil {
			return RepayEvent{}, err
		}
		liquidator, err := accountStore.FindById(*e.LiquidatorAccountId)
		if err != nil {
			return RepayEvent{}, err
		}
		cToken, err := tokenStore.FindById(*e.CollateralTokenId)
		if err != nil {
			return RepayEvent{}, err
		}
		dToken, err := tokenStore.FindById(*e.DebtTokenId)
		if err != nil {
			return RepayEvent{}, err
		}

		return LiquidationEvent{e.ID, e.Type, e.Txhash, e.Blocknumber, e.OccuredAt, borrower, liquidator, *e.AmountRepayed, *e.AmountSeized, cToken, dToken}, nil
	}
	return BorrowEvent{}, nil
}
