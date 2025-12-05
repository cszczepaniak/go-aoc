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

// Deprecated: WithInputFile shouldn't be used anymore. Main now caches inputs for you.
//
// TODO: let users turn off the caching maybe
func WithInputFile(path string) mainConfigOpt {
	return func(cfg mainConfig) mainConfig {
		cfg.inputPath = path
		return cfg
	}
}

type Solution[T any, In interface{ ~string | ~[]byte }] func(input In) T

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
func Main[T any, In interface{ ~string | ~[]byte }](
	year, day int,
	part1, part2 Solution[T, In],
	opts ...mainConfigOpt,
) error {
	submit := false
	clean := false
	flag.BoolVar(&submit, `submit`, false, `Whether or not to submit the answers from any provided solutions.`)
	flag.BoolVar(&clean, `clean`, false, `Whether or not to delete any cached input file for this year/day.`)
	flag.Parse()

	cachePath, err := cachedInputPath(year, day)
	if err != nil {
		return err
	}

	cfg := mainConfig{
		cl: http.DefaultClient,
	}
	for _, opt := range opts {
		cfg = opt(cfg)
	}

	if clean {
		err := os.RemoveAll(cachePath)
		if err != nil {
			return err
		}
	}

	input, err := GetAndCacheInput(context.Background(), cfg.cl, NewRequest(year, day).BuildGetInputRequest())
	if err != nil {
		return err
	}

	for _, s := range []struct {
		level    AnswerLevel
		solution Solution[T, In]
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

		res := s.solution(In(input))
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
