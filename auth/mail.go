package auth

// MailService specifies an API for sending emails via SMTP.
type MailService interface {

	// SendEmail sends an HTML email using given template and set of template variables.
	SendEmail(to string, subject string, templateName string, templateData map[string]interface{}) error
}
