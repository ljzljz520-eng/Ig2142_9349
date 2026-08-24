package app

import (
	"errors"
	"strings"
)

type Config struct {
	DatabasePath string
	Address      string
	Seed         bool
}

func DefaultConfig() Config { return Config{DatabasePath: "othello.db", Address: ":8080"} }

func (c Config) Validate() error {
	if strings.TrimSpace(c.DatabasePath) == "" {
		return errors.New("database path is required")
	}
	if strings.TrimSpace(c.Address) == "" {
		return errors.New("listen address is required")
	}
	return nil
}

func (c Config) WithDatabase(path string) Config {
	if strings.TrimSpace(path) != "" {
		c.DatabasePath = path
	}
	return c
}

func (c Config) WithAddress(address string) Config {
	if strings.TrimSpace(address) != "" {
		c.Address = address
	}
	return c
}
