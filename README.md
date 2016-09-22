# 2captcha
Little wrapper for 2captcha API

[![Build Status](https://travis-ci.org/LeoCBS/2captcha.svg?branch=master)](https://travis-ci.org/LeoCBS/2captcha)

## Go get

    go get github.com/LeoCBS/2captcha


## Upload image base64

    twocaptcha, _ := captcha.New("yourKey")
    captchaID, err := twocaptcha.UploadBase64Image("dHdvY2FwdGNoYQ==")

## Polling result

    iniAveragePollingTime := 1
    pollingTime := 1
    solution, err := twocaptcha.PollingCaptchaResponse(captchaID, iniAveragePollingTime, pollingTime)
    if err != nil {
	t.Errorf("unable to get captcha response: %s", err)
    }
