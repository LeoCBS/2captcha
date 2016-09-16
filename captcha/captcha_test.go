package captcha_test

import (
	"github.com/leocbs/2captcha/captcha"
	"github.com/leocbs/2captcha/captcha/httpmock"
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

func TestShouldValidateCaptchaIDPollingTime(t *testing.T) {
	twocaptcha, _ := captcha.New(testKey)
	_, err := twocaptcha.PollingCaptchaResponse("",1,1)
	if err == nil {
		t.Error("captchaID didn't validated")
	}
}


//test nil params
//test captcha solution
//test status code different 200