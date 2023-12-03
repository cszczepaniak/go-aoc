package aoc

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
)

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

func WithInputFile(path string) mainConfigOpt {
	return func(cfg mainConfig) mainConfig {
		cfg.inputPath = path
		return cfg
	}
}

type Solution[T any] func(input string) T

// Main runs the given solutions for the given puzzle year and day. Main looks for the `-submit` flag provided to the
// program to determine whether the answer should be submitted, or just printed to stdout. If submitted, the returned
// answer from each solution is formatted to a string using the `%v` verb.
//
// By default, the input will be downloaded from the Advent of Code website using the environment variable `AOC_SESSION`
// to authenticate. If an input file path is provided, that will be used instead. `AOC_SESSION` is also used for answer
// submission, if the `-submit` flag is set.
//
// A nil solution may be provided. Nil solutions will be skipped.
func Main[T any](
	year, day int,
	part1, part2 Solution[T],
	opts ...mainConfigOpt,
) error {
	submit := false
	flag.BoolVar(&submit, `submit`, false, `Whether or not to submit the answers from any provided solutions.`)
	flag.Parse()

	cfg := mainConfig{}
	for _, opt := range opts {
		cfg = opt(cfg)
	}

	var input string
	var err error

	switch {
	case cfg.inputPath != ``:
		bs, err := os.ReadFile(cfg.inputPath)
		if err != nil {
			return err
		}
		input = string(bs)
	case cfg.cl != nil:
		input, err = GetInputString(
			context.Background(),
			cfg.cl,
			NewRequest(year, day).BuildGetInputRequest(),
		)
	default:
		return errors.New(`no valid input source was provided`)
	}
	if err != nil {
		return err
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

		res := s.solution(input)
		fmt.Printf("Solution for part %v: %v\n", s.level, res)
		if !submit {
			continue
		}

		if cfg.cl != nil {
			fmt.Println("The submit flag was provided, but no HTTP client was configured.")
			continue
		}

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
			return err
		}
	}

	return nil
}
