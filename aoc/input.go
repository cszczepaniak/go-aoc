package aoc

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

type GetInputRequest struct {
	// Year is the year of the puzzle to fetch input for.
	Year int
	// Day is the day of the puzzle to fetch input for.
	Day int
	// SessionKey is the session key used to authenticate. If not set, it will be sourced from the AOC_SESSION
	// environment variable.
	SessionKey string
}

func (r GetInputRequest) url() string {
	return fmt.Sprintf(`https://adventofcode.com/%d/day/%d/input`, r.Year, r.Day)
}

func GetInputString(ctx context.Context, cl *http.Client, inputReq GetInputRequest) (string, error) {
	bs, err := GetInputBytes(ctx, cl, inputReq)
	if err != nil {
		return ``, err
	}

	return string(bs), nil
}

func GetInputBytes(ctx context.Context, cl *http.Client, inputReq GetInputRequest) ([]byte, error) {
	r, err := GetInput(ctx, cl, inputReq)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	return io.ReadAll(r)
}

func GetInput(ctx context.Context, cl *http.Client, inputReq GetInputRequest) (io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, inputReq.url(), nil)
	if err != nil {
		return nil, err
	}
	cookie, err := makeSessionCookie(inputReq.SessionKey)
	if err != nil {
		return nil, err
	}
	req.AddCookie(cookie)

	resp, err := cl.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf(`advent of code returned non-200 code: %d`, resp.StatusCode)
	}

	return resp.Body, nil
}

func makeSessionCookie(sessionKey string) (*http.Cookie, error) {
	if sessionKey == `` {
		var ok bool
		sessionKey, ok = os.LookupEnv(`AOC_SESSION`)
		if !ok {
			return nil, errors.New(`no session key found for authentication`)
		}
	}
	return &http.Cookie{
		Name:  `session`,
		Value: sessionKey,
	}, nil
}
