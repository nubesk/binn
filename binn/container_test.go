package binn

import (
	"time"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBottleGetter(t *testing.T) {
	d := time.Date(2022, 5, 29, 22, 24, 0, 0, time.UTC)
	b := NewBottle(
		"1c7a8201-cdf7-11ec-a9b3-0242ac110004",
		"This is a Test Message",
		&d,
	)

	d_ := time.Date(2022, 5, 29, 22, 24, 0, 0, time.UTC)
	assert.Equal(t, b.ID(), "1c7a8201-cdf7-11ec-a9b3-0242ac110004")
	assert.Equal(t, b.Message().Text, "This is a Test Message")
	assert.Equal(t, *b.ExpiredAt(), d_)
}
