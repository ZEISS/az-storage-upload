package main

import (
	"template/internal/cfg"

	"github.com/sethvargo/go-githubactions"
)

func main() {
	action := githubactions.New()

	cfg, err := cfg.NewFromInput(action)
	if err != nil {
		githubactions.Fatalf("error: %s", err)
	}

	_ = cfg
}
