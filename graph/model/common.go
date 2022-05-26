package model

import (
	"os"

	"log"

	"github.com/jmoiron/sqlx"
)

type Stores struct {
	Protocol         *ProtocolStore
	Chain            *ChainStore
	ProtocolInstance *ProtocolInstanceStore
	Account          *AccountStore
	Event            *EventStore
	Token            *TokenStore
}

func GenerateStores(db *sqlx.DB) Stores {

	protocolStore := NewProtocolStore(db)
	chainStore := NewChainStore(db)
	protocolInstanceStore := NewProtocolInstanceStore(db)
	accountStore := NewAccountStore(db)
	eventStore := NewEventStore(db)
	tokenStore := NewTokenStore(db)

	stores := Stores{
		Protocol:         protocolStore,
		Chain:            chainStore,
		ProtocolInstance: protocolInstanceStore,
		Account:          accountStore,
		Event:            eventStore,
		Token:            tokenStore,
	}

	return stores
}

func DbConnect() (*sqlx.DB, error) {
	url := os.Getenv("DATABASE_URL")
	if len(url) == 0 {
		log.Println("DATABASE_URL is not set")
	}
	return sqlx.Connect("pgx", url)
}
