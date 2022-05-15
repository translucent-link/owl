package model

import (
	"os"

	"github.com/jmoiron/sqlx"
)

type ChainStore struct {
	db *sqlx.DB
}

func NewChainStore() (*ChainStore, error) {
	db, err := sqlx.Connect("postgres", os.Getenv("DATABASE_URL"))
	return &ChainStore{db: db}, err
}

func (s *ChainStore) FindById(id int) (*Chain, error) {
	var chain Chain
	err := s.db.Get(&chain, "select * from chains where id=$1", id)
	return &chain, err
}

func (s *ChainStore) FindByName(name string) (*Chain, error) {
	var chain Chain
	err := s.db.Get(&chain, "select * from chains where name=$1", name)
	return &chain, err
}

func (s *ChainStore) CreateChain(input NewChain) (*Chain, error) {
	var insertedId int
	err := s.db.QueryRowx("insert into chains (name, rpcUrl, blockFetchSize) values ($1,$2,$3) returning id", input.Name, input.RPCURL, input.BlockFetchSize).Scan(&insertedId)
	if err != nil {
		return &Chain{}, err
	}
	return s.FindById(insertedId)
}

func (s *ChainStore) All() ([]*Chain, error) {
	chains := []*Chain{}
	err := s.db.Select(&chains, "SELECT * FROM chains ORDER BY name ASC")
	return chains, err
}
