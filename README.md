# 2captcha
Little wrapper for 2captcha API


## Upload image base64

    twocaptcha, _ := captcha.New("yourKey")
    captchaID, err := twocaptcha.UploadBase64Image("dHdvY2FwdGNoYQ==")