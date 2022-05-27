package model

import (
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type TokenStore struct {
	db *sqlx.DB
}

func NewTokenStore(db *sqlx.DB) *TokenStore {
	return &TokenStore{db: db}
}

func (s *TokenStore) FindById(id int) (*Token, error) {
	var token Token
	err := s.db.Get(&token, "select id, address, name, ticker, decimals from tokens where id=$1", id)
	return &token, err
}

func (s *TokenStore) FindByName(name string) (*Token, error) {
	var token Token
	err := s.db.Get(&token, "select id, address, name, ticker, decimals from tokens where name=$1", name)
	return &token, err
}

func (s *TokenStore) CreateToken(address string, name *string, ticker *string, chainId int) (*Token, error) {
	var insertedId int
	err := s.db.QueryRowx("insert into tokens (address, name, ticker, chainId) values ($1,$2,$3,$4) returning id", address, name, ticker, chainId).Scan(&insertedId)
	if err != nil {
		return &Token{}, err
	}
	return s.FindById(insertedId)
}

func (s *TokenStore) All() ([]*Token, error) {
	tokens := []*Token{}
	err := s.db.Select(&tokens, "SELECT id, address, name, ticker, decimals FROM tokens ORDER BY name ASC")
	return tokens, err
}

func (s *TokenStore) AllByChain(chainId int) ([]*Token, error) {
	tokens := []*Token{}
	err := s.db.Select(&tokens, "SELECT id, address, name, ticker, decimals FROM tokens where chainId=$1 ORDER BY name ASC", chainId)
	return tokens, err
}

func (s *TokenStore) FindOrCreateByAddress(address string, chainId int) (*Token, error) {
	var token Token
	err := s.db.Get(&token, "select id, address, name, ticker, decimals from tokens where address=$1", address)
	if err != nil {
		return s.CreateToken(address, nil, nil, chainId)
	}
	return &token, err
}

func (s *TokenStore) UpdateToken(id int, address string, name string, ticker string, decimals int) (*Token, error) {
	_, err := s.db.Exec("update tokens set address=$2, name=$3, ticker=$4, decimals=$5::int where id=$1", id, address, name, ticker, decimals)
	token, _ := s.FindById(id)
	return token, errors.Wrap(err, "Unable to update token")
}
