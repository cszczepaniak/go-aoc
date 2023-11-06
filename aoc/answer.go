package aoc

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type AnswerLevel int

const (
	AnswerPartOne AnswerLevel = 1
	AnswerPartTwo AnswerLevel = 2
)

var ErrWrongAnswer = errors.New(`incorrect answer`)

func SubmitAnswer(ctx context.Context, cl *http.Client, answerReq *submitAnswerRequest) error {
	req, err := answerReq.toHTTPRequest(ctx)
	if err != nil {
		return err
	}

	resp, err := cl.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf(`advent of code returned non-200 code: %d`, resp.StatusCode)
	}

	sc := bufio.NewScanner(resp.Body)
	for sc.Scan() {
		if strings.Contains(sc.Text(), `That's not the right answer;`) {
			return ErrWrongAnswer
		}
	}

	return nil
}
