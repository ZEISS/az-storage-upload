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
		"INPUT_NAME": "foo",
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
	assert.Equal(t, "foo", cfg.Name)
	assert.Equal(t, "", actionLog.String())
}
