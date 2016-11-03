package sendcloud

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type EmailSender struct {
	apiUser string
	apiKey  string
}

func NewEmailSender(apiUser string, apiKey string) *EmailSender {
	return &EmailSender{apiUser, apiKey}
}

type EmailData struct {
	From     string
	To       string
	Subject  string
	Html     string
	FromName string
	ReplyTo  string
}

type EmailSendResp struct {
	StatusCode int             `json:"statusCode"`
	Message    string          `json:"message"`
	Result     bool            `json:"result"`
	Info       json.RawMessage `json:"info"`
}

func (s *EmailSender) SendEmail(data *EmailData) error {
	host := "http://api.sendcloud.net/apiv2/mail/send"
	log.Printf(" SendEmail >>> %s \r\n", host)

	param := url.Values{}
	param.Set("apiUser", s.apiUser)
	param.Set("apiKey", s.apiKey)
	param.Set("from", data.From)
	param.Set("to", data.To)
	param.Set("subject", data.Subject)
	param.Set("html", data.Html)
	param.Set("fromName", data.FromName)
	param.Set("replyTo", data.ReplyTo)

	log.Printf(" param: %+v \r\n", param)

	resp, err := http.PostForm(host, param)
	if err != nil {
		log.Printf("邮件发送失败: %s \r\n", err.Error())
		return err
	}

	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("邮件发送失败: %s \r\n", err.Error())
		return err
	}
	log.Printf(" <<< resp: %s \r\n", string(bs))
	defer resp.Body.Close()

	var r EmailSendResp
	err = json.Unmarshal(bs, &r)
	if err != nil {
		log.Printf("邮件发送失败: %s \r\n", err.Error())
		return err
	}

	if r.StatusCode != 200 {
		log.Println("邮件发送失败: statusCode != 200")
		return errors.New("statusCode != 200")
	}

	if !r.Result {
		log.Println("邮件发送失败: !result")
		return errors.New("!result")
	}
	log.Println(" 邮件发送成功")

	return nil
}
