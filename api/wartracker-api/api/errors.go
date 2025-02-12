package api

import (
	"errors"
	"fmt"
)

var (
	ErrConfigFile       = errors.New("fatal error config file")
	ErrNoAliiances      = errors.New("no alliances found")
	ErrAliianceNotFound = errors.New("alliance not found")
	ErrVsDuelNotFound   = errors.New("vsduel not found")
	ErrNotImplemented   = errors.New("api not implemented")
)

func configError(err error) error {
	return fmt.Errorf("%w: %w", ErrConfigFile, err)
}
