package binn

import (
	"time"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	c := NewConfig(23, time.Duration(24 * time.Second),
		false, time.Duration(24 * time.Hour), false)

	assert.Equal(t, 23, c.Seed())
	assert.Equal(t, 24.0, c.DeliveryCycle().Seconds())
	assert.False(t, c.Validation())
	assert.Equal(t, 24.0, c.GenerateCycle().Hours())
	assert.False(t, c.Debug())
}

func TestDefaultConfig(t *testing.T) {
	c := DefaultConfig()

	assert.Equal(t, 42, c.Seed())
	assert.Equal(t, 15.0, c.DeliveryCycle().Minutes())
	assert.True(t, c.Validation())
	assert.Equal(t, 15.0, c.GenerateCycle().Minutes())
	assert.False(t, c.Debug())
}

func TestConfigSetter(t *testing.T) {
	c := DefaultConfig()

	c.SetSeed(5)
	c.SetDeliveryCycle(time.Duration(100) * time.Minute)
	c.SetGenerateCycle(time.Duration(33) * time.Minute)

	assert.Equal(t, 5, c.Seed())
	assert.Equal(t, 100.0, c.DeliveryCycle().Minutes())
	assert.Equal(t, 33.0, c.GenerateCycle().Minutes())
}

func TestConfigEnableBoolValue(t *testing.T) {
	c := NewConfig(42, time.Duration(1)*time.Minute,
		false, time.Duration(1)*time.Minute, false)
	c.EnableValidation()
	c.EnableDebug()

	assert.True(t, c.Validation())
	assert.True(t, c.Debug())
}

func TestConfigDisableBoolValue(t *testing.T) {
	c := NewConfig(42, time.Duration(1)*time.Minute,
		false, time.Duration(1)*time.Minute, false)
	c.DisableValidation()
	c.DisableDebug()

	assert.False(t, c.Validation())
	assert.False(t, c.Debug())
}
