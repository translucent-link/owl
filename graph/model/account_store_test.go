package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	setupDb()
}

func TestAllAccount(t *testing.T) {
	accountStore, _ := NewAccountStore()
	accounts, err := accountStore.All()
	assert.Nil(t, err)
	assert.Equal(t, 2, len(accounts))
}

func TestCreateAccount(t *testing.T) {
	accountStore, _ := NewAccountStore()
	account, err := accountStore.CreateAccount("0xB5DB0Eb39522427f292F4aeCA62B7886639BE8Dc")
	assert.Nil(t, err)
	assert.Equal(t, "0xB5DB0Eb39522427f292F4aeCA62B7886639BE8Dc", account.Address)
}

func TestFindAccountById(t *testing.T) {
	accountStore, _ := NewAccountStore()
	account, err := accountStore.FindById(1)
	assert.Nil(t, err)
	assert.Equal(t, "0xB5DB0Eb39522427f292F4aeCA62B7886639BE8Da", account.Address)
}

func TestFindAccountByAddress(t *testing.T) {
	accountStore, _ := NewAccountStore()
	account, err := accountStore.FindByAddress("0xB5DB0Eb39522427f292F4aeCA62B7886639BE8Da")
	assert.Nil(t, err)
	assert.Equal(t, 1, account.ID)
}
