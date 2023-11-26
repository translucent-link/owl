package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	setupDb()
}

func TestAllChains(t *testing.T) {
	db, _ := DbConnect()
	defer db.Close()
	stores := GenerateStores(db)

	chains, err := stores.Chain.All()
	assert.Nil(t, err)
	assert.Equal(t, 2, len(chains))
}

func TestCreateChain(t *testing.T) {
	db, _ := DbConnect()
	defer db.Close()
	stores := GenerateStores(db)

	chain, err := stores.Chain.CreateChain(NewChain{Name: "testnet"})
	assert.Nil(t, err)
	assert.Equal(t, "testnet", chain.Name)
	assert.NotEqual(t, 0, chain.ID)
}

func TestFindChainByName(t *testing.T) {
	db, _ := DbConnect()
	defer db.Close()
	stores := GenerateStores(db)

	chain, err := stores.Chain.FindByName("Avalanche")
	assert.Nil(t, err)
	assert.Equal(t, "Avalanche", chain.Name)
	assert.NotEqual(t, 0, chain.ID)
}

func TestFindChainByNameDoesntExist(t *testing.T) {
	db, _ := DbConnect()
	defer db.Close()
	stores := GenerateStores(db)

	_, err := stores.Chain.FindByName("Avalanche2")
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "sql: no rows in result set")
}

func TestFindChainById(t *testing.T) {
	db, _ := DbConnect()
	defer db.Close()
	stores := GenerateStores(db)

	chain, err := stores.Chain.FindById(1)
	assert.Nil(t, err)
	assert.Equal(t, "Avalanche", chain.Name)
	assert.NotEqual(t, 0, chain.ID)
}

func TestDeleteChainById(t *testing.T) {
	db, _ := DbConnect()
	defer db.Close()
	stores := GenerateStores(db)

	err := stores.Chain.DeleteByName("Sepolia")
	assert.Nil(t, err)
	_, err = stores.Chain.FindByName("Sepolia")
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "sql: no rows in result set")
}

func TestFetchAllChains(t *testing.T) {
	db, _ := DbConnect()
	defer db.Close()
	stores := GenerateStores(db)

	chains, err := stores.Chain.All()
	assert.Nil(t, err)
	assert.Equal(t, 2, len(chains))
}
