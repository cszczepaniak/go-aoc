package aoc

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

func cachedInputPath() (string, error) {
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

	return path, nil
}

type mainConfig struct {
	cl        *http.Client
	inputPath string
}

type mainConfigOpt func(cfg mainConfig) mainConfig

func WithHTTPClient(cl *http.Client) mainConfigOpt {
	return func(cfg mainConfig) mainConfig {
		cfg.cl = cl
		return cfg
	}
}

func WithDefaultHTTPClient() mainConfigOpt {
	return WithHTTPClient(http.DefaultClient)
}

// Deprecated: WithInputFile shouldn't be used anymore. Main now caches inputs for you.
//
// TODO: let users turn off the caching maybe
func WithInputFile(path string) mainConfigOpt {
	return func(cfg mainConfig) mainConfig {
		cfg.inputPath = path
		return cfg
	}
}

type Solution[T any] func(input string) T

// Main runs the given solutions for the given puzzle year and day. Main looks for the `-submit`
// flag provided to the program to determine whether the answer should be submitted, or just printed
// to stdout. The returned answer from each solution is formatted to a string using the `%v` verb.
//
// By default, the input will be downloaded from the Advent of Code website using the environment
// variable `AOC_SESSION` to authenticate. After the first download, the input will be saved for
// future runs.
//
// The given opts can be used to configure the HTTP client to use when downloading inputs and
// submitting answers. Other options are ignored.
//
// A nil solution may be provided. Nil solutions will be skipped.
func Main[T any](
	year, day int,
	part1, part2 Solution[T],
	opts ...mainConfigOpt,
) error {
	submit := false
	clean := false
	flag.BoolVar(&submit, `submit`, false, `Whether or not to submit the answers from any provided solutions.`)
	flag.BoolVar(&clean, `clean`, false, `Whether or not to delete any cached input file for this year/day.`)
	flag.Parse()

	cachePath, err := cachedInputPath()
	if err != nil {
		return err
	}

	cfg := mainConfig{
		cl: http.DefaultClient,
	}
	for _, opt := range opts {
		cfg = opt(cfg)
	}

	var input []byte
	dirPath := filepath.Join(cachePath, strconv.Itoa(year), strconv.Itoa(day))
	fullPath := filepath.Join(dirPath, "input.txt")
	found := false
	if clean {
		err := os.RemoveAll(fullPath)
		if err != nil {
			return err
		}
	} else {
		input, err = os.ReadFile(fullPath)
		found = !os.IsNotExist(err)
		if err != nil && found {
			return err
		}
	}

	if !found {
		input, err = GetInputBytes(
			context.Background(),
			cfg.cl,
			NewRequest(year, day).BuildGetInputRequest(),
		)
		if err != nil {
			return err
		}

		if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
			return err
		}

		if err := os.WriteFile(fullPath, input, os.ModePerm); err != nil {
			return err
		}
	}

	for _, s := range []struct {
		level    AnswerLevel
		solution Solution[T]
	}{{
		level:    AnswerPartOne,
		solution: part1,
	}, {
		level:    AnswerPartTwo,
		solution: part2,
	}} {
		if s.solution == nil {
			fmt.Printf("No solution provided for part %d\n", s.level)
			continue
		}

		res := s.solution(string(input))
		fmt.Printf("Solution for part %v: %v\n", s.level, res)
		if !submit {
			continue
		}

		if cfg.cl == nil {
			fmt.Println("The submit flag was provided, but no HTTP client was configured.")
			continue
		}

		fmt.Printf("Submitting answer for part %v... ", s.level)

		err = SubmitAnswer(
			context.Background(),
			cfg.cl,
			NewRequest(
				year,
				day,
			).BuildSubmitAnswerRequest(
				s.level,
				fmt.Sprintf("%v", res),
			),
		)
		if err != nil {
			if errors.Is(err, ErrWrongAnswer) {
				fmt.Println(`Sorry, wrong answer!`)
				continue
			}
			return err
		}
		fmt.Println(`Correct answer!`)
	}

	return nil
}
