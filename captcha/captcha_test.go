package captcha_test

import (
	"github.com/leocbs/2captcha/captcha"
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"
)

const testKey = "testKey"

type Responder func(*http.Request) (*http.Response, error)

type MockTransport struct {
	responders  map[string]Responder
	noResponder Responder
}

var defaultTransport = NewMockTransport()

func NewMockTransport() *MockTransport {
	return &MockTransport{make(map[string]Responder), nil}
}

func (m *MockTransport) registerResponder(method, url string, responder Responder) {
	m.responders[method+" "+url] = responder
}

func newRespBodyFromString(body string) io.ReadCloser {
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

//end

func setUp() {
	http.DefaultTransport = defaultTransport
}

func TestShouldValidateEmptyKey(t *testing.T) {
	_, err := captcha.New("")
	if err == nil {
		t.Error("new captcha don't valid empty key")
	}
}

func TestShouldValidateEmptyBase64(t *testing.T) {
	twocaptcha, _ := captcha.New(testKey)
	_, err := twocaptcha.UploadBase64Image("")
	if err == nil {
		t.Error("new captcha don't valid empty base64 image")
	}
}

func TestShouldUploadBase64Image(t *testing.T) {
	setUp()
	defaultTransport.registerResponder("POST", "http://2captcha.com/in.php",
		func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				Status:     strconv.Itoa(500),
				StatusCode: 200,
				Body:       newRespBodyFromString("OK|captchaID"),
				Header:     http.Header{},
			}, nil

		})

	twocaptcha, _ := captcha.New(testKey)
	_, err := twocaptcha.UploadBase64Image("dHdvY2FwdGNoYQ==")
	if err != nil {
		t.Errorf("upload response not OK: %s", err)
	}
}
