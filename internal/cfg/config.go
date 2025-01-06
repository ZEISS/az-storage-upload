package cfg

import (
	"errors"

	"github.com/sethvargo/go-githubactions"
	"github.com/zeiss/pkg/utilx"
)

var (
	// ErrMissingPath returns an error when the path is missing.
	ErrMissingPath = errors.New("missing path")
	// ErrMissingAccountURL returns an error when the account URL is missing.
	ErrMissingAccountURL = errors.New("missing account URL")
	// ErrMissingContainerName returns an error when the container name is missing.
	ErrMissingContainerName = errors.New("missing container name")
)

// Config ...
type Config struct {
	Path          string
	AccountURL    string
	ContainerName string
}

// NewFromInput ...
func NewFromInput(action *githubactions.Action) (*Config, error) {
	cfg := new(Config)

	cfg.Path = action.GetInput("path")
	if utilx.Empty(cfg.Path) {
		return nil, ErrMissingPath
	}

	cfg.AccountURL = action.GetInput("account_url")
	if utilx.Empty(cfg.AccountURL) {
		return nil, ErrMissingAccountURL
	}

	cfg.ContainerName = action.GetInput("container_name")
	if utilx.Empty(cfg.ContainerName) {
		return nil, ErrMissingContainerName
	}

	return cfg, nil
}
