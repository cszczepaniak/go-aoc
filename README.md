# go-aoc
Advent of Code helpers in Go. I use these to download input and submit answers programmatically.

As Eric Wastl (creator of AoC) points out in his HTML source:

> Please be careful with automated requests; I'm not a massive company, and I can only take so much traffic.  Please be 
considerate so that everyone gets to play.

Use these helpers responsibly!

## Examples

### Getting puzzle input

```go
r, err := GetInput(
    context.Background(),
    http.DefaultClient,
    NewRequest(2015, 1).BuildGetInputRequest(),
)
if err != nil {
    return err
}
defer r.Close()

// Do stuff with `r` (which is an io.Reader)

bs, err := GetInputBytes(
    context.Background(),
    http.DefaultClient,
    NewRequest(2015, 1).BuildGetInputRequest(),
)
if err != nil {
    return err
}

// Do stuff with `bs` (which is a slice of bytes)

str, err := GetInputString(
    context.Background(),
    http.DefaultClient,
    NewRequest(2015, 1).BuildGetInputRequest(),
)
if err != nil {
    return err
}

// Do stuff with `str` (which is a string)
```

### Submitting answers

```go
err := SubmitAnswer(
    context.Background(),
    http.DefaultClient,
    NewRequest(2015, 7).BuildSubmitAnswerRequest(AnswerPartOne, `123`),
)
if err != nil {
    if errors.Is(err, aoc.ErrWrongAnswer) {
        // Successful request, but your answer was wrong.
    }
    return err
}
```

### Authentication

By default, a request creating using `NewRequest(year, day)` reads a session key from the `AOC_SESSION` environment 
variable. This can be customized.

```go
// Use a different environment variable to load the session key.
req := NewRequest(2015, 7).WithSessionKeyFromEnv(`MY_ENV_VAR`)

// Set the session key explicitly.
req := NewRequest(2015, 7).WithSessionKey(`my_secret_session_key`)
```