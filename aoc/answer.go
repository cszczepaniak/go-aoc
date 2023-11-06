package aoc

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type AnswerLevel int

const (
	AnswerPartOne AnswerLevel = 1
	AnswerPartTwo AnswerLevel = 2
)

type SubmitAnswerRequest struct {
	// Year is the year of the puzzle to submit an answer for.
	Year int
	// Day is the day of the puzzle to submit an answer for.
	Day int
	// Answer is the answer to submit.
	Answer string
	// Level is the level of answer being submitted. Valid options are AnswerPartOne and AnswerPartTwo.
	Level AnswerLevel
	// SessionKey is the session key used to authenticate. If not set, it will be sourced from the AOC_SESSION
	// environment variable.
	SessionKey string
}

func (r SubmitAnswerRequest) url() string {
	return fmt.Sprintf(`https://adventofcode.com/%d/day/%d/answer`, r.Year, r.Day)
}

func (r SubmitAnswerRequest) requestBody() io.Reader {
	vals := url.Values{}
	vals.Add(`level`, strconv.Itoa(int(r.Level)))
	vals.Add(`answer`, r.Answer)

	return strings.NewReader(vals.Encode())
}

var errWrongAnswer = errors.New(`incorrect answer`)

func SubmitAnswer(ctx context.Context, cl *http.Client, answerReq SubmitAnswerRequest) error {
	req, err := http.NewRequest(http.MethodPost, answerReq.url(), answerReq.requestBody())
	if err != nil {
		return err
	}

	cookie, err := makeSessionCookie(answerReq.SessionKey)
	if err != nil {
		return err
	}

	req.AddCookie(cookie)
	req.Header.Add(`Content-Type`, `application/x-www-form-urlencoded`)

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
			return errWrongAnswer
		}
	}

	return nil
}
