package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	setupDb()
}

func TestAllProtocolInstances(t *testing.T) {
	db, _ := DbConnect()
	defer db.Close()
	stores := GenerateStores(db)

	protocolInstances, err := stores.ProtocolInstance.All()
	assert.Nil(t, err)
	assert.Equal(t, 3, len(protocolInstances))
}

func TestCreateProtocolInstance(t *testing.T) {
	db, _ := DbConnect()
	defer db.Close()
	stores := GenerateStores(db)

	_, _ = stores.Chain.CreateChain(NewChain{Name: "Testnet3"})

	protocolInstance, err := stores.ProtocolInstance.CreateProtocolInstance(NewProtocolInstance{Protocol: "Aave", Chain: "Testnet3", ContractAddress: "0x1234", FirstBlockToRead: 10000})
	assert.Nil(t, err)
	assert.Equal(t, "Testnet3", protocolInstance.Chain.Name)
	assert.Equal(t, "Aave", protocolInstance.Protocol.Name)
	assert.Equal(t, "0x1234", protocolInstance.ContractAddress)
	assert.Equal(t, 10000, protocolInstance.FirstBlockToRead)
	assert.Equal(t, 0, protocolInstance.LastBlockRead)
	assert.NotEqual(t, 0, protocolInstance.ID)
}

func TestFindProtocolInstanceByChainAndProtocol(t *testing.T) {
	db, _ := DbConnect()
	defer db.Close()
	stores := GenerateStores(db)

	protocolInstance, err := stores.ProtocolInstance.FindByProtocolIdAndChainId(1, 1)
	assert.Nil(t, err)
	assert.Equal(t, "0xB5DB0Eb39522427f292F4aeCA62B7886639BE8De", protocolInstance.ContractAddress)
	assert.NotEqual(t, 0, protocolInstance.ID)
}

func TestFindProtocolInstanceByChainAndProtocolDoesntExist(t *testing.T) {
	db, _ := DbConnect()
	defer db.Close()
	stores := GenerateStores(db)

	_, err := stores.ProtocolInstance.FindByProtocolIdAndChainId(100, 100)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "sql: no rows in result set")
}

func TestFindProtocolInstanceById(t *testing.T) {
	db, _ := DbConnect()
	defer db.Close()
	stores := GenerateStores(db)

	protocolInstance, err := stores.ProtocolInstance.FindById(1)
	assert.Nil(t, err)
	assert.Equal(t, "0xB5DB0Eb39522427f292F4aeCA62B7886639BE8De", protocolInstance.ContractAddress)
	assert.NotEqual(t, 0, protocolInstance.ID)
}

func TestDeleteProtocolInstanceByProtocolAndChain(t *testing.T) {
	db, _ := DbConnect()
	defer db.Close()
	stores := GenerateStores(db)

	err := stores.ProtocolInstance.DeleteByProtocolAndChain("Aave", "Arbitrum")
	assert.Nil(t, err)
	_, err = stores.ProtocolInstance.FindById(3)
	assert.NotNil(t, err)
	assert.Equal(t, "Unable to find protocol instance with id 3: sql: no rows in result set", err.Error())
}
