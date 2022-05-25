package model

import (
	"os"

	"github.com/jmoiron/sqlx"
)

func Stores() (*ChainStore, *ProtocolStore, *ProtocolInstanceStore, error) {
	protocolStore, err := NewProtocolStore()
	if err != nil {
		return &ChainStore{}, &ProtocolStore{}, &ProtocolInstanceStore{}, err
	}
	chainStore, err := NewChainStore()
	if err != nil {
		return &ChainStore{}, &ProtocolStore{}, &ProtocolInstanceStore{}, err
	}
	protocolInstanceStore, err := NewProtocolInstanceStore()
	if err != nil {
		return &ChainStore{}, &ProtocolStore{}, &ProtocolInstanceStore{}, err
	}

	return chainStore, protocolStore, protocolInstanceStore, nil
}

func DbConnect() (*sqlx.DB, error) {
	return sqlx.Connect("pgx", os.Getenv("DATABASE_URL"))
}
