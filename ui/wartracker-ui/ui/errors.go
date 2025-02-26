package ui

import (
	"errors"
	"fmt"
)

var (
	ErrConfigFile       = errors.New("fatal error config file")
	ErrNoAliiances      = errors.New("no alliances found")
	ErrAliianceNotFound = errors.New("alliance not found")
)

func configError(err error) error {
	return fmt.Errorf("%w: %w", ErrConfigFile, err)
}
