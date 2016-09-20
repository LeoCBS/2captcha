# 2captcha
Little wrapper for 2captcha API


## Upload image base64

    twocaptcha, _ := captcha.New("yourKey")
    captchaID, err := twocaptcha.UploadBase64Image("dHdvY2FwdGNoYQ==")

## Polling result

    twocaptcha, _ := captcha.New("yourKey")
    captchaID, err := twocaptcha.UploadBase64Image("dHdvY2FwdGNoYQ==")
    iniAveragePollingTime := 1
    pollingTime := 1
    solution, err := twocaptcha.PollingCaptchaResponse(captchaID, iniAveragePollingTime, pollingTime)
    if err != nil {
	t.Errorf("unable to get captcha response: %s", err)
    }
