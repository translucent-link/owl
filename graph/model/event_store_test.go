package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	setupDb()
}

func TestAllAccountEvents(t *testing.T) {
	eventStore, _ := NewEventStore()
	events, err := eventStore.AllByAccount(1)

	assert.Nil(t, err)
	assert.Equal(t, 3, len(events))
	assert.Equal(t, "Borrow", events[0].Type.String())
	assert.Equal(t, "Repay", events[1].Type.String())
	assert.Equal(t, "Liquidation", events[2].Type.String())
}
