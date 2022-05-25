package model

import (
	"github.com/jmoiron/sqlx"
)

type AccountStore struct {
	db *sqlx.DB
}

func NewAccountStore() (*AccountStore, error) {
	db, err := DbConnect()
	return &AccountStore{db: db}, err
}

func (s *AccountStore) FindById(id int) (*Account, error) {
	var account Account
	err := s.db.Get(&account, "select * from accounts where id=$1", id)
	return &account, err
}

func (s *AccountStore) FindByAddress(address string) (*Account, error) {
	var account Account
	err := s.db.Get(&account, "select * from accounts where address=$1", address)
	return &account, err
}

func (s *AccountStore) CreateAccount(address string) (*Account, error) {
	var insertedId int
	err := s.db.QueryRowx("insert into accounts (address) values ($1) returning id", address).Scan(&insertedId)
	if err != nil {
		return &Account{}, err
	}
	return s.FindById(insertedId)
}

func (s *AccountStore) All() ([]*Account, error) {
	accounts := []*Account{}
	err := s.db.Select(&accounts, "SELECT * FROM accounts ORDER BY address ASC")
	return accounts, err
}

func (s *AccountStore) FindOrCreateByAddress(address string) (*Account, error) {
	account, err := s.FindByAddress(address)
	if err != nil {
		return s.CreateAccount(address)
	}
	return account, err
}
