package captcha

import "errors"

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


func UploadImageFormData(file []byte) (string, error){
	return "", nil
}


// polling 2captcha response page until captcha is ready.
// initSleep represent 2captcha average time to solve captcha, don't make senses polling
// response before average time
func PollingCaptchaResponse(captchaId string, initSleep int, pollingTime int) (string, error){
	return "", nil
}


