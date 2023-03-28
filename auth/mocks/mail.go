package mocks

import (
	"github.com/datapace/datapace/auth"
)

var _ auth.MailService = (*mailServiceMock)(nil)

type mailServiceMock struct {
	smtpIdentity    string
	smtpHost        string
	smtpURL         string
	smtpUser        string
	smtpPassword    string
	smtpFrom        string
	frontendURL     string
	passRecoveryTpl string
}

func NewMailService() auth.MailService {
	return &mailServiceMock{}
}

func (ms *mailServiceMock) SendRecoveryEmail(to string, subject string, templateData map[string]interface{}) error {

	return nil
}
