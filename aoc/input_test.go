package aoc

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetInput_Integration(t *testing.T) {
	sess, ok := os.LookupEnv(`AOC_SESSION`)
	if !ok {
		t.Skip(`AOC_SESSION variable is not set`)
	}

	r, err := GetInput(context.Background(), http.DefaultClient, GetInputRequest{
		Year: 2015,
		Day:  1,
	})
	require.NoError(t, err)
	defer r.Close()

	bs1, err := io.ReadAll(r)
	require.NoError(t, err)
	assert.NotEmpty(t, bs1)

	bs2, err := GetInputBytes(context.Background(), http.DefaultClient, GetInputRequest{
		Year:       2015,
		Day:        1,
		SessionKey: sess,
	})
	require.NoError(t, err)

	assert.True(t, bytes.Equal(bs1, bs2), `expected bytes to be equal`)
}
