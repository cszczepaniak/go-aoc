package aoc

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

func GetAndCacheInput(ctx context.Context, cl *http.Client, inputReq *getInputRequest) ([]byte, error) {
	if cl == nil {
		cl = http.DefaultClient
	}

	cachePath, err := cachedInputPath(inputReq.year, inputReq.day)
	if err != nil {
		return nil, err
	}

	dirPath := filepath.Dir(cachePath)

	bs, err := os.ReadFile(cachePath)
	if os.IsNotExist(err) {
		bs, err = GetInputBytes(ctx, cl, inputReq)
		if err != nil {
			return nil, err
		}

		err = os.MkdirAll(dirPath, os.ModePerm)
		if err != nil {
			return nil, err
		}

		err = os.WriteFile(cachePath, bs, os.ModePerm)
		if err != nil {
			return nil, err
		}

		return bs, nil
	} else if err != nil {
		return nil, err
	}

	return bs, nil
}

func cachedInputPath(year, day int) (string, error) {
	home, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	path := filepath.Join(home, "aoc")

	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			return "", err
		}
	} else if err != nil {
		return "", err
	}

	return filepath.Join(path, strconv.Itoa(year), strconv.Itoa(day), "input.txt"), nil
}

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
