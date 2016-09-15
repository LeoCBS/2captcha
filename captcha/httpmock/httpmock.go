package httpmock

import (
	"io"
	"net/http"
	"strings"
)

type Responder func(*http.Request) (*http.Response, error)

type MockTransport struct {
	responders  map[string]Responder
	noResponder Responder
}

var DefaultTransport = NewMockTransport()

func NewMockTransport() *MockTransport {
	return &MockTransport{make(map[string]Responder), nil}
}

func (m *MockTransport) RegisterResponder(method, url string, responder Responder) {
	m.responders[method+" "+url] = responder
}

func NewRespBodyFromString(body string) io.ReadCloser {
	return &dummyReadCloser{strings.NewReader(body)}
}
func (m *MockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	url := req.URL.String()

	// try and get a responder that matches the method and URL
	responder := m.responderForKey(req.Method + " " + url)

	// if we weren't able to find a responder and the URL contains a querystring
	// then we strip off the querystring and try again.
	if responder == nil && strings.Contains(url, "?") {
		responder = m.responderForKey(req.Method + " " + strings.Split(url, "?")[0])
	}

	// if we found a responder, call it
	if responder != nil {
		return responder(req)
	}

	return nil, nil
}
func (m *MockTransport) responderForKey(key string) Responder {
	for k, r := range m.responders {
		if k != key {
			continue
		}
		return r
	}
	return nil
}

// dummycloser
type dummyReadCloser struct {
	body io.ReadSeeker
}

func (d *dummyReadCloser) Read(p []byte) (n int, err error) {
	n, err = d.body.Read(p)
	if err == io.EOF {
		d.body.Seek(0, 0)
	}
	return n, err
}

func (d *dummyReadCloser) Close() error {
	return nil
}