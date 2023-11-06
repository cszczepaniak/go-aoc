package aoc

import (
	"context"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubmitAnswer_Integration(t *testing.T) {
	_, ok := os.LookupEnv(`AOC_SESSION`)
	if !ok {
		t.Skip(`AOC_SESSION variable is not set`)
	}

	err := SubmitAnswer(
		context.Background(),
		http.DefaultClient,
		NewRequest(2015, 7).BuildSubmitAnswerRequest(AnswerPartOne, `123`),
	)
	assert.ErrorIs(t, err, ErrWrongAnswer)
}
