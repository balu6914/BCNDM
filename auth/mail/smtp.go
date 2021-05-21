package mail

import (
	"bytes"
	"fmt"
	"github.com/datapace/datapace/auth"
	"html/template"
	"net/smtp"
)

var _ auth.MailService = (*mailService)(nil)

type mailService struct {
	smtpIdentity    string
	smtpHost        string
	smtpURL         string
	smtpUser        string
	smtpPassword    string
	smtpFrom        string
	frontendURL     string
	passRecoveryTpl string
}

func New(smtpIdentity string, smtpURL string, smtpHost string, smtpUser string, smtpPassword string, smtpFrom string, frontendURL string, passRecoveryTpl string) auth.MailService {
	return &mailService{
		smtpIdentity,
		smtpHost,
		smtpURL,
		smtpUser,
		smtpPassword,
		smtpFrom,
		frontendURL,
		passRecoveryTpl,
	}
}

func (ms *mailService) SendRecoveryEmail(to string, subject string, templateData map[string]interface{}) error {
	auth := smtp.PlainAuth(ms.smtpIdentity, ms.smtpUser, ms.smtpPassword, ms.smtpHost)
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

	smtpErr := smtp.SendMail(ms.smtpURL, auth, ms.smtpFrom, []string{to}, body.Bytes())

	if smtpErr != nil {
		return smtpErr
	}

	return nil
}
