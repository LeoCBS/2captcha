package captcha

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strings"
	"time"
)

const (
	inputUrl    = "http://2captcha.com/in.php"
	responseUrl = "http://2captcha.com/res.php"
	OK          = "OK"
	notReady    = "CAPCHA_NOT_READY"
)

type Captcha struct {
	key string
}

func New(key string) (*Captcha, error) {
	if key == "" {
		return nil, errors.New("key should not be empty")
	}
	return &Captcha{
		key: key,
	}, nil
}

//upload one base64 imagem to twocaptcha to be resolverd
//return captcha id or one error
func (captcha *Captcha) UploadBase64Image(base64 string) (string, error) {
	if base64 == "" {
		return "", errors.New("base64 should be not empty")
	}
	bf, contentType, err := captcha.createForm(base64)
	if err != nil {
		return "", errors.New("failed to create form")
	}

	req, err := http.NewRequest("POST", inputUrl, bf)
	if err != nil {
		return "", errors.New("failed to create request/post")
	}
	req.Header.Set("Content-Type", contentType)
	body, err := perfomRequest(req)
	if err != nil {
		return "", err
	}
	captchaID, err := getValueOK(body)
	if err != nil {
		return "", err
	}
	return captchaID, nil
}

func getValueOK(body string) (string, error) {
	if strings.Contains(body, "OK|") {
		return strings.Split(body, "|")[1], nil
	}
	return "", errors.New(body)
}

func perfomRequest(request *http.Request) (string, error) {
	client := &http.Client{}
	resp, err := client.Do(request)
	defer func() {
		if resp.Body != nil {
			resp.Body.Close()
		}
	}()

	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("response status code different 200")
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

type formCreator struct {
	err error
}

func (captcha *Captcha) createForm(image string) (*bytes.Buffer, string, error) {
	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)
	defer writer.Close()

	formCreator := &formCreator{}
	formCreator.createFormField("key", captcha.key, writer)
	formCreator.createFormField("body", image, writer)
	formCreator.createFormField("method", "base64", writer)
	if formCreator.err != nil {
		return nil, "", formCreator.err
	}

	return &buffer, writer.FormDataContentType(), nil
}

func (fc *formCreator) createFormField(fieldName string, fieldValue string, writer *multipart.Writer) {
	if fc.err != nil {
		return
	}
	fw, err := writer.CreateFormField(fieldName)
	if err != nil {
		fc.err = errors.New(fmt.Sprintf("failed to create field %s ", fieldName))
		return
	}
	if _, err := fw.Write([]byte(fieldValue)); err != nil {
		fc.err = errors.New(fmt.Sprintf("failed to set %s value", fieldName))
		return
	}
}

// polling 2captcha response page until captcha is ready.
// initAverageSleep represent 2captcha average time to solve captcha, don't makes senses polling
// response before average time
func (captcha *Captcha) PollingCaptchaResponse(captchaID string, initAverageSleep time.Duration, pollingTime time.Duration) (string, error) {
	validator := &pollingValidator{}
	validator.validatePollingParams(captchaID, initAverageSleep, pollingTime)
	if validator.err != nil {
		return "", validator.err
	}

	time.Sleep(initAverageSleep)
	body, err := captcha.getResponse(captchaID, pollingTime)
	if err != nil {
		return "", err
	}
	solution, err := getValueOK(body)
	if err != nil {
		return "", err
	}
	return solution, nil
}

type pollingValidator struct {
	err error
}

func (validator *pollingValidator) validatePollingParams(captchaID string,
	initAverageSleep time.Duration,
	pollingTime time.Duration) {
	if validator.err != nil {
		return
	}
	if captchaID == "" {
		validator.err = errors.New("CaptchaID should not be empty")
		return
	}
	if initAverageSleep == 0 {
		validator.err = errors.New("initAverageSleep should not be zero")
		return
	}
	if pollingTime == 0 {
		validator.err = errors.New("pollingTime should not be zero")
		return
	}

}

func (captcha *Captcha) getResponse(captchaID string, pollingTime time.Duration) (string, error) {
	req, _ := http.NewRequest("GET", responseUrl, nil)

	q := req.URL.Query()
	q.Add("key", captcha.key)
	q.Add("action", "get")
	q.Add("id", captchaID)
	req.URL.RawQuery = q.Encode()
	body, err := perfomRequest(req)
	if err != nil {
		return "", err
	}
	if body == notReady {
		time.Sleep(pollingTime)
		return captcha.getResponse(captchaID, pollingTime)
	}
	return body, nil
}
