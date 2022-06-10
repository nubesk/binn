package binn

import (
	"time"
)


type Config struct {
	seed          int
	deliveryCycle time.Duration
	validation    bool
	generateCycle time.Duration
	debug         bool
}

func NewConfig(s int, d time.Duration, v bool, g time.Duration, ed bool) *Config {
	return &Config{
		seed:          s,
		deliveryCycle: d,
		validation:    v,
		generateCycle: g,
		debug:         ed,
	}
}

func DefaultConfig() *Config {
	return &Config{
		seed:          42,
		deliveryCycle: time.Duration(15 * time.Minute),
		validation:    true,
		generateCycle: time.Duration(15 * time.Minute),
		debug:         false,
	}
}

func (c *Config) Seed() int {
	return c.seed
}

func (c *Config) SetSeed(s int) {
	c.seed = s
}

func (c *Config) DeliveryCycle() time.Duration {
	return c.deliveryCycle
}

func (c *Config) SetDeliveryCycle(d time.Duration) {
	c.deliveryCycle = d
}

func (c *Config) Validation() bool {
	return c.validation
}

func (c *Config) EnableValidation() {
	c.validation = true
}

func (c *Config) DisableValidation() {
	c.validation = false
}

func (c *Config) GenerateCycle() time.Duration {
	return c.generateCycle
}

func (c *Config) SetGenerateCycle(g time.Duration) {
	c.generateCycle = g
}

func (c *Config) Debug() bool {
	return c.debug
}

func (c *Config) EnableDebug() {
	c.debug = true
}

func (c *Config) DisableDebug() {
	c.debug = false
}
