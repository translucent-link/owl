package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	setupDb()
}

func TestAllAccount(t *testing.T) {
	db, _ := DbConnect()
	defer db.Close()
	stores := GenerateStores(db)

	accounts, err := stores.Account.All()
	assert.Nil(t, err)
	assert.Equal(t, 2, len(accounts))
}

func TestCreateAccount(t *testing.T) {
	db, _ := DbConnect()
	defer db.Close()
	stores := GenerateStores(db)

	account, err := stores.Account.CreateAccount("0xB5DB0Eb39522427f292F4aeCA62B7886639BE8Dc")
	assert.Nil(t, err)
	assert.Equal(t, "0xB5DB0Eb39522427f292F4aeCA62B7886639BE8Dc", account.Address)
}

func TestFindAccountById(t *testing.T) {
	db, _ := DbConnect()
	defer db.Close()
	stores := GenerateStores(db)

	account, err := stores.Account.FindById(1)
	assert.Nil(t, err)
	assert.Equal(t, "0xB5DB0Eb39522427f292F4aeCA62B7886639BE8Da", account.Address)
}

func TestFindAccountByAddress(t *testing.T) {
	db, _ := DbConnect()
	defer db.Close()
	stores := GenerateStores(db)

	account, err := stores.Account.FindByAddress("0xB5DB0Eb39522427f292F4aeCA62B7886639BE8Da")
	assert.Nil(t, err)
	assert.Equal(t, 1, account.ID)
}
