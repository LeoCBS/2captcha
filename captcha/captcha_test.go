package captcha_test

import (
	"github.com/LeoCBS/2captcha/captcha"
	"github.com/LeoCBS/2captcha/captcha/httpmock"
	"net/http"
	"strconv"
	"testing"
)

const testKey = "testKey"

func setUp() {
	http.DefaultTransport = httpmock.DefaultTransport
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
	httpmock.DefaultTransport.RegisterResponder("POST", "http://2captcha.com/in.php",
		func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				Status:     strconv.Itoa(200),
				StatusCode: 200,
				Body:       httpmock.NewRespBodyFromString("OK|captchaID"),
				Header:     http.Header{},
			}, nil

		})

	twocaptcha, _ := captcha.New(testKey)
	captchaID, err := twocaptcha.UploadBase64Image("dHdvY2FwdGNoYQ==")
	if err != nil {
		t.Errorf("upload response not OK: %s", err)
	}
	if captchaID != "captchaID" {
		t.Error("upload image don't return captcha id correctly")
	}
}

func TestShouldValidateCaptchaIDPollingResponse(t *testing.T) {
	twocaptcha, _ := captcha.New(testKey)
	_, err := twocaptcha.PollingCaptchaResponse("", 1, 1)
	if err == nil {
		t.Error("captchaID didn't validated")
	}
}

func TestShouldValidateInitAverageSleepPollingResponse(t *testing.T) {
	twocaptcha, _ := captcha.New(testKey)
	_, err := twocaptcha.PollingCaptchaResponse("captchaID", 0, 1)
	if err == nil {
		t.Error("sleep init time didn't validated")
	}
}

func TestShouldValidatePollingTimePollingResponse(t *testing.T) {
	twocaptcha, _ := captcha.New(testKey)
	_, err := twocaptcha.PollingCaptchaResponse("captchaID", 1, 0)
	if err == nil {
		t.Error("polling time didn't validated")
	}
}

func TestShouldPollingCaptchaResponse(t *testing.T) {
	setUp()
	httpmock.DefaultTransport.RegisterResponder("GET", "http://2captcha.com/res.php",
		func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				Status:     strconv.Itoa(200),
				StatusCode: 200,
				Body:       httpmock.NewRespBodyFromString("OK|captchaBroken"),
				Header:     http.Header{},
			}, nil

		})

	twocaptcha, _ := captcha.New(testKey)
	captchaID, err := twocaptcha.UploadBase64Image("dHdvY2FwdGNoYQ==")
	if err != nil {
		t.Errorf("upload response not OK: %s", err)
	}
	if captchaID != "captchaID" {
		t.Error("upload image don't return captcha id correctly")
	}
	solution, err := twocaptcha.PollingCaptchaResponse(captchaID, 1, 1)
	if err != nil {
		t.Errorf("unable to get captcha response: %s", err)
	}
	if solution != "captchaBroken" {
		t.Error("wrong captcha solution")
	}
}

//test status code different 200
