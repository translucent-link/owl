package model

import (
	"github.com/jmoiron/sqlx"
)

type TokenStore struct {
	db *sqlx.DB
}

func NewTokenStore() (*TokenStore, error) {
	db, err := DbConnect()
	return &TokenStore{db: db}, err
}

func (s *TokenStore) FindById(id int) (*Token, error) {
	var token Token
	err := s.db.Get(&token, "select id, address, name, ticker from tokens where id=$1", id)
	return &token, err
}

func (s *TokenStore) FindByName(name string) (*Token, error) {
	var token Token
	err := s.db.Get(&token, "select id, address, name, ticker from tokens where name=$1", name)
	return &token, err
}

func (s *TokenStore) CreateToken(address string, name *string, ticker *string) (*Token, error) {
	var insertedId int
	err := s.db.QueryRowx("insert into tokens (address, name, ticker) values ($1,$2,$3) returning id", address, name, ticker).Scan(&insertedId)
	if err != nil {
		return &Token{}, err
	}
	return s.FindById(insertedId)
}

func (s *TokenStore) All() ([]*Token, error) {
	tokens := []*Token{}
	err := s.db.Select(&tokens, "SELECT id, address, name, ticker FROM tokens ORDER BY name ASC")
	return tokens, err
}

func (s *TokenStore) AllByChain(chainId int) ([]*Token, error) {
	tokens := []*Token{}
	err := s.db.Select(&tokens, "SELECT id, address, name, ticker FROM tokens where chainId=$1 ORDER BY name ASC", chainId)
	return tokens, err
}

func (s *TokenStore) FindOrCreateByAddress(address string) (*Token, error) {
	var token Token
	err := s.db.Get(&token, "select id, address, name, ticker from tokens where address=$1", address)
	if err != nil {
		return s.CreateToken(address, nil, nil)
	}
	return &token, err
}
