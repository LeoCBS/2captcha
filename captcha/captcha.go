package captcha

import (
	"strings"
	"bytes"
	"errors"
	"net/http"
	"io/ioutil"
	"mime/multipart"
	"fmt"
	)

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
	_, err = getCaptchaID(body)
	if err != nil{
		return "", err
	}
	return body, nil
}

func getCaptchaID(body string) (string, error){
	if strings.Contains(body, "OK|"){
		return strings.Split(body,"|")[1], nil
	}
	return "", errors.New(body)
}

func perfomRequest(request *http.Request) (string, error){
	client := &http.Client{}
	resp, err := client.Do(request)
	defer resp.Body.Close()

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
// initSleep represent 2captcha average time to solve captcha, don't make senses polling
// response before average time
func PollingCaptchaResponse(captchaId string, initSleep int, pollingTime int) (string, error){
	return "", nil
}


