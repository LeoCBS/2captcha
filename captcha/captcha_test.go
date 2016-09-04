package captcha_test

import (
	"testing"

    "github.com/leocbs/2captcha/captcha"
)

const testKey = "testKey"

func TestShouldValidateEmptyKey(t *testing.T) {
    _,err := captcha.New("")
    if err == nil{
        t.Error("new captcha don't valid empty key")
    }
}

func TestShouldValidateEmptyBase64(t *testing.T) {
	twocaptcha,_ := captcha.New(testKey)
	_, err := twocaptcha.UploadBase64Image("")
	if err == nil{
        t.Error("new captcha don't valid empty base64 image")
    } 
}