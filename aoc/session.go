package aoc

import "os"

type namedEnvSessionKeyProvider struct {
	envName string
}

func (p namedEnvSessionKeyProvider) getSessionKey() string {
	return os.Getenv(p.envName)
}

type constSessionKeyProvider string

func (p constSessionKeyProvider) getSessionKey() string {
	return string(p)
}
