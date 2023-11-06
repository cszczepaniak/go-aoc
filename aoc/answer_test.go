package aoc

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/cszczepaniak/go-aoc/aoc/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSubmitAnswer_Integration(t *testing.T) {
	_, ok := os.LookupEnv(`AOC_SESSION`)
	if !ok || testing.Short() {
		t.Skip(`AOC_SESSION variable is not set`)
	}

	err := SubmitAnswer(
		context.Background(),
		http.DefaultClient,
		NewRequest(2015, 7).BuildSubmitAnswerRequest(AnswerPartOne, `123`),
	)
	assert.ErrorIs(t, err, ErrWrongAnswer)
}

func TestSubmitAnswer(t *testing.T) {
	h := testutils.NewHTTPTestHandler(t)
	s := httptest.NewServer(h)
	t.Cleanup(s.Close)

	err := SubmitAnswer(
		context.Background(),
		http.DefaultClient,
		newRequest(s.URL, 2015, 7).WithSessionKey(`foo`).BuildSubmitAnswerRequest(AnswerPartOne, `123`),
	)
	require.NoError(t, err)

	req := h.LastRequest()
	assert.Equal(t, `/2015/day/7/answer`, req.Req.RequestURI)
	assert.Equal(t, http.MethodPost, req.Req.Method)
	assert.Equal(t, `application/x-www-form-urlencoded`, req.Req.Header.Get(`Content-Type`))
	cookie, err := req.Req.Cookie(`session`)
	require.NoError(t, err)
	assert.Equal(t, `foo`, cookie.Value)

	h.RespondWith(http.StatusOK, "liasbdakljshdbas\naksjlhdkjhdasThat's not the right answer;kjsahdklasjd\n\n\n\n")
	err = SubmitAnswer(
		context.Background(),
		http.DefaultClient,
		newRequest(s.URL, 2015, 7).WithSessionKey(`foo`).BuildSubmitAnswerRequest(AnswerPartOne, `123`),
	)
	assert.ErrorIs(t, err, ErrWrongAnswer)

	h.RespondWith(http.StatusForbidden, `foobar`)
	err = SubmitAnswer(
		context.Background(),
		http.DefaultClient,
		newRequest(s.URL, 2015, 7).WithSessionKey(`foo`).BuildSubmitAnswerRequest(AnswerPartOne, `123`),
	)
	assert.EqualError(t, err, `advent of code returned non-200 code: 403`)
}
