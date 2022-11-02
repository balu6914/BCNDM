package mail

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"strconv"

	"github.com/datapace/datapace/auth"
	"gopkg.in/gomail.v2"
)

var _ auth.MailService = (*mailService)(nil)

type mailService struct {
	smtpIdentity    string
	smtpHost        string
	smtpPort        string
	smtpURL         string
	smtpUser        string
	smtpPassword    string
	smtpFrom        string
	frontendURL     string
	passRecoveryTpl string
}

func New(smtpIdentity string, smtpURL string, smtpHost string, smtpPort string, smtpUser string, smtpPassword string, smtpFrom string, frontendURL string, passRecoveryTpl string) auth.MailService {
	return &mailService{
		smtpIdentity,
		smtpHost,
		smtpPort,
		smtpURL,
		smtpUser,
		smtpPassword,
		smtpFrom,
		frontendURL,
		passRecoveryTpl,
	}
}

func (ms *mailService) SendRecoveryEmail(to string, subject string, templateData map[string]interface{}) error {
	port, err := strconv.Atoi(ms.smtpPort)
	if err != nil {
		return err
	}
	dial := gomail.NewDialer(ms.smtpHost, port, ms.smtpUser, ms.smtpPassword)
	dial.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	header := "To:" + to + "\r\n" + "From:" + ms.smtpFrom + "\r\n" + "Subject:" + subject
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";"
	mailTemplate, parsingErr := template.ParseFiles(ms.passRecoveryTpl)
	if parsingErr != nil {
		return parsingErr
	}

	body := bytes.Buffer{}
	body.Write([]byte(fmt.Sprintf("%s\n%s\n\n", header, mime)))
	templateData["FrontendURL"] = ms.frontendURL
	templateErr := mailTemplate.Execute(&body, templateData)
	if templateErr != nil {
		return templateErr
	}

	m := gomail.NewMessage()
	m.SetHeader("From", ms.smtpFrom)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body.String())

	if err := dial.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
