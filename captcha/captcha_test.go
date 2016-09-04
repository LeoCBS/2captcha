package captcha_test

import (
	"testing"
	"io"
	"net/http"

    "github.com/leocbs/2captcha/captcha"
)

const testKey = "testKey"

//start http server to receive and test two captcha behavior
func setUp(){
	http.HandleFunc("/", mockIinputCaptcha)
	http.ListenAndServe(":8000", nil)
}

func mockIinputCaptcha(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello world!")
}

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