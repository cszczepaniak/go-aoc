package testutils

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

type Request struct {
	Req  *http.Request
	Body string
}

type HTTPTestServer struct {
	t            testing.TB
	respondWith  int
	responseBody string

	ReceivedRequests []Request
}

func NewHTTPTestHandler(t testing.TB) *HTTPTestServer {
	return &HTTPTestServer{
		t:           t,
		respondWith: http.StatusOK,
	}
}

func (s *HTTPTestServer) RespondWith(statusCode int, body string) {
	s.respondWith = statusCode
	s.responseBody = body
}

func (s *HTTPTestServer) Reset() {
	s.ReceivedRequests = nil
}

func (s *HTTPTestServer) LastRequest() Request {
	require.NotEmpty(s.t, s.ReceivedRequests, `no requests found`)

	return s.ReceivedRequests[len(s.ReceivedRequests)-1]
}

func (s *HTTPTestServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	bs, err := io.ReadAll(r.Body)
	require.NoError(s.t, err)

	s.ReceivedRequests = append(s.ReceivedRequests, Request{
		Req:  r,
		Body: string(bs),
	})

	w.WriteHeader(s.respondWith)
	if s.responseBody != `` {
		w.Write([]byte(s.responseBody))
	}
}
