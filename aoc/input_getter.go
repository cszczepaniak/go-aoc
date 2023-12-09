package aoc

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
)

type inputGetter struct {
	cachePath string
	cl        *http.Client
}

func newInputGetter(cachePath string, cl *http.Client) *inputGetter {
	return &inputGetter{
		cachePath: cachePath,
		cl:        cl,
	}
}

func (ig *inputGetter) GetInputString(ctx context.Context, inputReq *getInputRequest) (string, error) {
	bs, err := ig.GetInputBytes(ctx, inputReq)
	if err != nil {
		return ``, err
	}

	return string(bs), nil
}

func (ig *inputGetter) GetInputBytes(ctx context.Context, inputReq *getInputRequest) ([]byte, error) {
	if ig.cachePath != `` {
		bs, err := ig.getCachedInput(inputReq)
		if err != nil {
			return nil, fmt.Errorf(`failed to get cached input: %w`, err)
		}
		return bs, nil
	}

	req, err := inputReq.toHTTPRequest(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := ig.cl.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf(`advent of code returned non-200 code: %d`, resp.StatusCode)
	}

	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if ig.cachePath != `` {
		err := ig.putCachedInput(inputReq, bs)
		if err != nil {
			return nil, fmt.Errorf(`failed to cache input: %w`, err)
		}
	}

	return bs, nil
}

func (ig *inputGetter) getCachedInput(req *getInputRequest) ([]byte, error) {
	p := path.Join(
		ig.cachePath,
		strconv.Itoa(req.year),
		strconv.Itoa(req.day),
		`input.txt`,
	)
	return os.ReadFile(p)
}

func (ig *inputGetter) putCachedInput(req *getInputRequest, data []byte) error {
	p, err := ig.getOrCreateCachePath(req)
	if err != nil {
		return err
	}
	return os.WriteFile(p, data, 0o755)
}

func (ig *inputGetter) getOrCreateCachePath(req *getInputRequest) (string, error) {
	p := path.Join(
		ig.cachePath,
		strconv.Itoa(req.year),
		strconv.Itoa(req.day),
	)
	err := os.MkdirAll(p, os.ModePerm)
	if err != nil {
		return ``, err
	}
	return path.Join(p, `input.txt`), nil
}
