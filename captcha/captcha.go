package captcha

import "errors"

const (
	inputUrl  = "http://2captcha.com/in.php"
	responseUrl = "http://2captcha.com/res.php"
	OK             = "OK"
	notReady       = "CAPCHA_NOT_READY"
	reportedOK     = "OK_REPORT_RECORDED"
)

type Captcha struct{
	key string
}

func New(key string) (*Captcha, error){
	if key == ""{
		return nil, errors.New("key should not be empty")
	}
	return &Captcha{
		key: key,
	}, nil
}


//upload one base64 imagem to twocaptcha to be resolverd
//return captcha id or one error
func (captcha *Captcha) UploadBase64Image(base64 string) (string, error){
	if base64 == ""{
		return "", errors.New("base64 should be not empty")
	}

	return "", nil
}

// polling 2captcha response page until captcha is ready.
// initSleep represent 2captcha average time to solve captcha, don't make senses polling
// response before average time
func PollingCaptchaResponse(captchaId string, initSleep int, pollingTime int) (string, error){
	return "", nil
}


