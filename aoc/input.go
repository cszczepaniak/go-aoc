package aoc

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

func GetInputString(ctx context.Context, cl *http.Client, inputReq *getInputRequest) (string, error) {
	bs, err := GetInputBytes(ctx, cl, inputReq)
	if err != nil {
		return ``, err
	}

	return string(bs), nil
}

func GetInputBytes(ctx context.Context, cl *http.Client, inputReq *getInputRequest) ([]byte, error) {
	r, err := GetInput(ctx, cl, inputReq)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	return io.ReadAll(r)
}

func GetInput(ctx context.Context, cl *http.Client, inputReq *getInputRequest) (io.ReadCloser, error) {
	req, err := inputReq.toHTTPRequest(ctx)
	if err != nil {
		return nil, err
	}

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
