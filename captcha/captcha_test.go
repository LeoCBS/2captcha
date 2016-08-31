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

func TestShouldValidateEmptyFile(t *testing.T) {
    //_,_ := captcha.New(testKey)
	
    
}