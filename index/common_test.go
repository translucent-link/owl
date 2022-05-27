package index

import (
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func init() {
}

func TestGetEstimatedTimestamp(t *testing.T) {

	startTime := time.Date(2022, 5, 27, 2, 37, 29, 0, time.UTC)
	endTime := time.Date(2022, 5, 27, 8, 56, 35, 0, time.UTC)
	// expected := time.Date(2022, 5, 27, 3, 54, 58, 0, time.UTC)
	acceptable := time.Date(2022, 5, 27, 2, 57, 24, 584756898, time.UTC)
	estimate := GetEstimatedTimestamp(startTime, endTime, big.NewInt(14851700), big.NewInt(14853222), big.NewInt(14851780))

	assert.Equal(t, acceptable, estimate)
}
