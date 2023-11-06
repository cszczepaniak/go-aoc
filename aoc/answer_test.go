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

	err := SubmitAnswer(context.Background(), http.DefaultClient, SubmitAnswerRequest{
		Year:   2015,
		Day:    7,
		Answer: `123`,
		Level:  AnswerPartOne,
	})
	assert.ErrorIs(t, err, errWrongAnswer)
}
