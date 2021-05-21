package auth

// MailService specifies an API for sending emails via SMTP.
type MailService interface {

	// SendRecoveryEmail sends an password recovery email using a set of template variables.
	SendRecoveryEmail(to string, subject string, templateData map[string]interface{}) error
}
