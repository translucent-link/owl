package model

import (
	"errors"

	"github.com/jmoiron/sqlx"
)

type ProtocolStore struct {
	db *sqlx.DB
}

func NewProtocolStore(db *sqlx.DB) *ProtocolStore {
	return &ProtocolStore{db: db}
}

func (s *ProtocolStore) FindById(id int) (*Protocol, error) {
	var protocol Protocol
	err := s.db.Get(&protocol, "select * from protocols where id=$1", id)
	return &protocol, err
}

func (s *ProtocolStore) FindByName(name string) (*Protocol, error) {
	var protocol Protocol
	err := s.db.Get(&protocol, "select * from protocols where name=$1", name)
	return &protocol, err
}

func (s *ProtocolStore) CreateProtocol(input NewProtocol) (*Protocol, error) {
	var insertedId int
	err := s.db.QueryRowx("insert into protocols (name, abi) values ($1, $2) returning id", input.Name, input.Abi).Scan(&insertedId)
	if err != nil {
		return &Protocol{}, err
	}
	return s.FindById(insertedId)
}

func (s *ProtocolStore) All() ([]*Protocol, error) {
	protocols := []*Protocol{}
	err := s.db.Select(&protocols, "SELECT * FROM protocols ORDER BY name ASC")
	return protocols, err
}

func (s *ProtocolStore) AllByChain(chainId int) ([]*Protocol, error) {
	protocols := []*Protocol{}
	err := s.db.Select(&protocols, "SELECT p.* FROM protocols p JOIN protocol_instances pi ON p.id=pi.protocolId WHERE pi.chainId=$1 ORDER BY name ASC", chainId)
	return protocols, err
}

func (s *ProtocolStore) AddEventDefn(protocolName string, topicName string, topicHashHex string, abiSignature string) (*EventDefn, error) {
	protocol, err := s.FindByName(protocolName)
	if err != nil {
		return &EventDefn{}, errors.New("Unable to find protocol")
	}
	var insertedId int
	err = s.db.QueryRowx("insert into event_definitions (protocolId, topicName, topicHashHex, abiSignature) values ($1, $2, $3, $4) returning id", protocol.ID, topicName, topicHashHex, abiSignature).Scan(&insertedId)
	if err != nil {
		return &EventDefn{}, err
	}
	return s.FindEventById(insertedId)
}

func (s *ProtocolStore) FindEventById(id int) (*EventDefn, error) {
	var eventDefn EventDefn
	err := s.db.Get(&eventDefn, "select id, topicName, topicHashHex, abiSignature from event_definitions where id=$1", id)
	return &eventDefn, err
}

func (s *ProtocolStore) AllEventsByProtocol(id int) ([]*EventDefn, error) {
	events := []*EventDefn{}
	err := s.db.Select(&events, "select id, topicName, topicHashHex, abiSignature from event_definitions where protocolId=$1", id)
	return events, err
}

func (s *ProtocolStore) AllEventsByProtocolAndTopicName(id int, topicName string) ([]*EventDefn, error) {
	events := []*EventDefn{}
	err := s.db.Select(&events, "select id, topicName, topicHashHex, abiSignature from event_definitions where protocolId=$1 and topicName=$2", id, topicName)
	return events, err
}
