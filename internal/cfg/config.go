package cfg

import (
	"github.com/sethvargo/go-githubactions"
)

// Config ...
type Config struct {
	Name string
}

// NewFromInput ...
func NewFromInput(action *githubactions.Action) (*Config, error) {
	cfg := new(Config)
	cfg.Name = action.GetInput("name")

	return cfg, nil
}
