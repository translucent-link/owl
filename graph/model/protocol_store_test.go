package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	setupDb()
}

func TestAllProtocols(t *testing.T) {
	db, _ := DbConnect()
	defer db.Close()
	stores := GenerateStores(db)

	protocols, err := stores.Protocol.All()
	assert.Nil(t, err)
	assert.Equal(t, 3, len(protocols))
}

func TestCreateProtocol(t *testing.T) {
	db, _ := DbConnect()
	defer db.Close()
	stores := GenerateStores(db)

	protocol, err := stores.Protocol.CreateProtocol(NewProtocol{Name: "compound", Abi: "someABI"})
	assert.Nil(t, err)
	assert.Equal(t, "compound", protocol.Name)
	assert.Equal(t, "someABI", protocol.Abi)
	assert.NotEqual(t, 0, protocol.ID)
}

func TestFindProtocolByName(t *testing.T) {
	db, _ := DbConnect()
	defer db.Close()
	stores := GenerateStores(db)

	protocol, err := stores.Protocol.FindByName("Aave")
	assert.Nil(t, err)
	assert.Equal(t, "Aave", protocol.Name)
	assert.NotEqual(t, "", protocol.Abi)
	assert.NotEqual(t, 0, protocol.ID)
}

func TestFindProtocolByNameDoesntExist(t *testing.T) {
	db, _ := DbConnect()
	defer db.Close()
	stores := GenerateStores(db)

	_, err := stores.Protocol.FindByName("Aave2")
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "sql: no rows in result set")
}

func TestFindProtocolById(t *testing.T) {
	db, _ := DbConnect()
	defer db.Close()
	stores := GenerateStores(db)

	protocol, err := stores.Protocol.FindById(1)
	assert.Nil(t, err)
	assert.Equal(t, "Aave", protocol.Name)
	assert.NotEqual(t, 0, protocol.ID)
}

func TestDeleteProtocolByName(t *testing.T) {
	db, _ := DbConnect()
	defer db.Close()
	stores := GenerateStores(db)

	err := stores.Protocol.DeleteByName("AaveToBeDeleted")
	assert.Nil(t, err)
	_, err = stores.Protocol.FindByName("AaveToBeDeleted")
	assert.NotNil(t, err)
	assert.Equal(t, "sql: no rows in result set", err.Error())
}
