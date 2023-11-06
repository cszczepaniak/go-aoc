package aoc

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type baseRequest struct {
	// Year is the year of the puzzle to submit an answer for.
	year int
	// Day is the day of the puzzle to submit an answer for.
	day int
	// SessionKey is the session key used to authenticate. If not set, it will be sourced from the AOC_SESSION
	// environment variable.
	sessionKeyProvider interface{ getSessionKey() string }
}

// NewRequest returns a new Advent of Code request. By default, the request will read the session key from the
// environment variable `AOC_SESSION`.
func NewRequest(year, day int) *baseRequest {
	return &baseRequest{
		year: year,
		day:  day,
		sessionKeyProvider: namedEnvSessionKeyProvider{
			envName: `AOC_SESSION`,
		},
	}
}

// WithSessionKeyFromEnv sets the environment variable name the request will read the session key from.
func (br *baseRequest) WithSessionKeyFromEnv(envName string) *baseRequest {
	br.sessionKeyProvider = namedEnvSessionKeyProvider{
		envName: envName,
	}
	return br
}

// WithSessionKey sets the session key of the request to the provided value.
func (br *baseRequest) WithSessionKey(sessionKey string) *baseRequest {
	br.sessionKeyProvider = constSessionKeyProvider(sessionKey)
	return br
}

func (br *baseRequest) baseURL() string {
	return fmt.Sprintf(`https://adventofcode.com/%d/day/%d`, br.year, br.day)
}

func (br *baseRequest) addCookie(req *http.Request) error {
	sessionKey := br.sessionKeyProvider.getSessionKey()
	if sessionKey == `` {
		return errors.New(`no session key found for authentication`)
	}

	req.AddCookie(&http.Cookie{
		Name:  `session`,
		Value: sessionKey,
	})
	return nil
}

func (br *baseRequest) BuildSubmitAnswerRequest(level AnswerLevel, answer string) *submitAnswerRequest {
	return &submitAnswerRequest{
		baseRequest: br,
		answer:      answer,
		level:       level,
	}
}

func (br *baseRequest) BuildGetInputRequest() *getInputRequest {
	return &getInputRequest{
		baseRequest: br,
	}
}

type submitAnswerRequest struct {
	*baseRequest
	// Answer is the answer to submit.
	answer string
	// Level is the level of answer being submitted. Valid options are AnswerPartOne and AnswerPartTwo.
	level AnswerLevel
}

func (r *submitAnswerRequest) toHTTPRequest(ctx context.Context) (*http.Request, error) {
	u := r.baseRequest.baseURL() + `/answer`

	vals := url.Values{}
	vals.Add(`level`, strconv.Itoa(int(r.level)))
	vals.Add(`answer`, r.answer)

	body := strings.NewReader(vals.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, body)
	if err != nil {
		return nil, err
	}
	err = r.addCookie(req)
	if err != nil {
		return nil, err
	}

	req.Header.Add(`Content-Type`, `application/x-www-form-urlencoded`)
	return req, nil
}

type getInputRequest struct {
	*baseRequest
}

func (r *getInputRequest) toHTTPRequest(ctx context.Context) (*http.Request, error) {
	u := r.baseRequest.baseURL() + `/input`

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, nil)
	if err != nil {
		return nil, err
	}
	err = r.addCookie(req)
	if err != nil {
		return nil, err
	}

	return req, nil
}
