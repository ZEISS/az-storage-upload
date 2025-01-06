package cfg

import (
	"bytes"
	"testing"

	"github.com/sethvargo/go-githubactions"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFromInputs(t *testing.T) {
	actionLog := bytes.NewBuffer(nil)

	envMap := map[string]string{
		"INPUT_PATH":           "/var/lib/data",
		"INPUT_ACCOUNT_URL":    "https://example.com",
		"INPUT_CONTAINER_NAME": "container",
	}

	getenv := func(key string) string {
		return envMap[key]
	}

	action := githubactions.New(
		githubactions.WithWriter(actionLog),
		githubactions.WithGetenv(getenv),
	)

	cfg, err := NewFromInput(action)
	require.NoError(t, err)

	assert.NotNil(t, cfg)
	assert.Equal(t, "/var/lib/data", cfg.Path)
	assert.Equal(t, "https://example.com", cfg.AccountURL)
	assert.Equal(t, "container", cfg.ContainerName)
	assert.Equal(t, "", actionLog.String())
}
