package aoc

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/cszczepaniak/go-aoc/aoc/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetInput_Integration(t *testing.T) {
	sess, ok := os.LookupEnv(`AOC_SESSION`)
	if !ok || testing.Short() {
		t.Skip(`AOC_SESSION variable is not set`)
	}

	r, err := GetInput(
		context.Background(),
		http.DefaultClient,
		NewRequest(2015, 1).BuildGetInputRequest(),
	)
	require.NoError(t, err)
	defer r.Close()

	bs1, err := io.ReadAll(r)
	require.NoError(t, err)
	assert.NotEmpty(t, bs1)

	bs2, err := GetInputBytes(
		context.Background(),
		http.DefaultClient,
		NewRequest(2015, 1).WithSessionKey(sess).BuildGetInputRequest(),
	)
	require.NoError(t, err)

	assert.True(t, bytes.Equal(bs1, bs2), `expected bytes to be equal`)
}

func TestGetInputString(t *testing.T) {
	h := testutils.NewHTTPTestHandler(t)
	s := httptest.NewServer(h)
	t.Cleanup(s.Close)

	h.RespondWith(http.StatusOK, `your puzzle input`)
	str, err := GetInputString(
		context.Background(),
		http.DefaultClient,
		newRequest(s.URL, 2015, 7).WithSessionKey(`foo`).BuildGetInputRequest(),
	)
	require.NoError(t, err)
	assert.Equal(t, `your puzzle input`, str)

	req := h.LastRequest()
	assert.Equal(t, `/2015/day/7/input`, req.Req.RequestURI)
	assert.Equal(t, http.MethodGet, req.Req.Method)
	cookie, err := req.Req.Cookie(`session`)
	require.NoError(t, err)
	assert.Equal(t, `foo`, cookie.Value)
}

func TestGetInputBytes(t *testing.T) {
	h := testutils.NewHTTPTestHandler(t)
	s := httptest.NewServer(h)
	t.Cleanup(s.Close)

	h.RespondWith(http.StatusOK, `your puzzle input`)
	bs, err := GetInputBytes(
		context.Background(),
		http.DefaultClient,
		newRequest(s.URL, 2015, 7).WithSessionKey(`foo`).BuildGetInputRequest(),
	)
	require.NoError(t, err)
	assert.True(t, bytes.Equal([]byte(`your puzzle input`), bs), `bytes were not equal`)

	req := h.LastRequest()
	assert.Equal(t, `/2015/day/7/input`, req.Req.RequestURI)
	assert.Equal(t, http.MethodGet, req.Req.Method)
	cookie, err := req.Req.Cookie(`session`)
	require.NoError(t, err)
	assert.Equal(t, `foo`, cookie.Value)
}

func TestGetInput(t *testing.T) {
	h := testutils.NewHTTPTestHandler(t)
	s := httptest.NewServer(h)
	t.Cleanup(s.Close)

	h.RespondWith(http.StatusOK, `your puzzle input`)
	r, err := GetInput(
		context.Background(),
		http.DefaultClient,
		newRequest(s.URL, 2015, 7).WithSessionKey(`foo`).BuildGetInputRequest(),
	)
	require.NoError(t, err)

	bs, err := io.ReadAll(r)
	require.NoError(t, err)
	assert.True(t, bytes.Equal([]byte(`your puzzle input`), bs), `bytes were not equal`)

	req := h.LastRequest()
	assert.Equal(t, `/2015/day/7/input`, req.Req.RequestURI)
	assert.Equal(t, http.MethodGet, req.Req.Method)
	cookie, err := req.Req.Cookie(`session`)
	require.NoError(t, err)
	assert.Equal(t, `foo`, cookie.Value)
}

func TestGetInputString_NotOK(t *testing.T) {
	h := testutils.NewHTTPTestHandler(t)
	s := httptest.NewServer(h)
	t.Cleanup(s.Close)

	h.RespondWith(http.StatusForbidden, `foobar`)
	_, err := GetInputString(
		context.Background(),
		http.DefaultClient,
		newRequest(s.URL, 2015, 7).WithSessionKey(`foo`).BuildGetInputRequest(),
	)
	assert.EqualError(t, err, `advent of code returned non-200 code: 403`)
}
