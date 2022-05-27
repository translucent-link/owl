package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	setupDb()
}

func TestAllTokens(t *testing.T) {
	db, _ := DbConnect()
	defer db.Close()
	stores := GenerateStores(db)

	tokens, err := stores.Token.All()
	assert.Nil(t, err)
	assert.Equal(t, 2, len(tokens))
}

func TestCreateToken(t *testing.T) {
	db, _ := DbConnect()
	defer db.Close()
	stores := GenerateStores(db)

	name := "Polygon Matic"
	ticker := "MATIC"
	token, err := stores.Token.CreateToken("0xB5DB0Eb39522427f292F4aeCA62B7886639BE8Db", &name, &ticker, 1)
	assert.Nil(t, err)
	assert.NotEqual(t, 0, token.ID)
	assert.Equal(t, "0xB5DB0Eb39522427f292F4aeCA62B7886639BE8Db", token.Address)
	assert.Equal(t, "Polygon Matic", *(token.Name))
	assert.Equal(t, "MATIC", *(token.Ticker))
}

func TestFindTokenById(t *testing.T) {
	db, _ := DbConnect()
	defer db.Close()
	stores := GenerateStores(db)

	token, err := stores.Token.FindById(1)
	assert.Nil(t, err)
	assert.Equal(t, 1, token.ID)
	assert.Equal(t, "0xB5DB0Eb39522427f292F4aeCA62B7886639BE8Dc", token.Address)
	assert.Equal(t, "Polygon", *(token.Name))
	assert.Equal(t, "MATIC", *(token.Ticker))
}

func TestFindTokenByName(t *testing.T) {
	db, _ := DbConnect()
	defer db.Close()
	stores := GenerateStores(db)

	token, err := stores.Token.FindByName("Polygon")
	assert.Nil(t, err)
	assert.Equal(t, 1, token.ID)
	assert.Equal(t, "0xB5DB0Eb39522427f292F4aeCA62B7886639BE8Dc", token.Address)
	assert.Equal(t, "Polygon", *(token.Name))
	assert.Equal(t, "MATIC", *(token.Ticker))
}
