package model

import (
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type ProtocolInstanceStore struct {
	db            *sqlx.DB
	protocolStore *ProtocolStore
	chainStore    *ChainStore
}

func NewProtocolInstanceStore() (*ProtocolInstanceStore, error) {
	db, err := DbConnect()
	if err != nil {
		return &ProtocolInstanceStore{}, err
	}
	ps, err := NewProtocolStore()
	if err != nil {
		return &ProtocolInstanceStore{}, err
	}
	cs, err := NewChainStore()
	return &ProtocolInstanceStore{db: db, protocolStore: ps, chainStore: cs}, err
}

func (s *ProtocolInstanceStore) FindById(id int) (*ProtocolInstance, error) {
	var protocolInstance ProtocolInstance
	err := s.db.Get(&protocolInstance, "select id, contractAddress, firstBlockToRead, lastBlockRead from protocol_instances where id=$1", id)
	return &protocolInstance, err
}

func (s *ProtocolInstanceStore) FindByProtocolIdAndChainId(protocolId int, chainId int) (*ProtocolInstance, error) {
	fmt.Printf("p %d c %d", protocolId, chainId)
	var protocolInstance ProtocolInstance
	err := s.db.Get(&protocolInstance, "select id, contractAddress, firstBlockToRead, lastBlockRead from protocol_instances where protocolId=$1 and chainId=$2", protocolId, chainId)
	return &protocolInstance, err
}

func (s *ProtocolInstanceStore) FindProtocolById(id int) (*Protocol, error) {
	var protocolId int
	var protocol Protocol
	err := s.db.QueryRowx("SELECT protocolId FROM protocol_instances WHERE id=$1", id).Scan(&protocolId)
	if err != nil {
		return &protocol, err
	}
	return s.protocolStore.FindById(protocolId)
}

func (s *ProtocolInstanceStore) FindChainById(id int) (*Chain, error) {
	var chainId int
	var chain Chain
	err := s.db.QueryRowx("SELECT chainId FROM protocol_instances WHERE id=$1", id).Scan(&chainId)
	if err != nil {
		return &chain, err
	}
	return s.chainStore.FindById(chainId)
}

func (s *ProtocolInstanceStore) CreateProtocolInstance(input NewProtocolInstance) (*ProtocolInstance, error) {
	var insertedId int

	protocol, err := s.protocolStore.FindByName(input.Protocol)
	if err != nil {
		return &ProtocolInstance{}, errors.New("Unable to find protocol")
	}
	chain, err := s.chainStore.FindByName(input.Chain)
	if err != nil {
		return &ProtocolInstance{}, errors.New("Unable to find chain")
	}

	err = s.db.QueryRowx("insert into protocol_instances (chainId, protocolId, contractAddress, firstBlockToRead) values ($1, $2, $3, $4::int) returning id", chain.ID, protocol.ID, input.ContractAddress, input.FirstBlockToRead).Scan(&insertedId)
	if err != nil {
		return &ProtocolInstance{}, err
	}
	fmt.Printf("Returning PI id: %d", insertedId)
	return s.FindById(insertedId)
}

func (s *ProtocolInstanceStore) All() ([]*ProtocolInstance, error) {
	protocolInstances := []*ProtocolInstance{}
	err := s.db.Select(&protocolInstances, "SELECT id, contractAddress FROM protocol_instances ORDER BY id ASC")
	return protocolInstances, err
}

func (s *ProtocolInstanceStore) UpdateLastBlockRead(protocolInstanceId int, lastBlock uint) error {
	result, err := s.db.Exec("update protocol_instances set lastBlockRead=$2::int where id=$1", protocolInstanceId, lastBlock)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New(fmt.Sprintf("Failed to update lastBlockRead=%d on protocolInstance %d", lastBlock, protocolInstanceId))
	}
	return nil
}
