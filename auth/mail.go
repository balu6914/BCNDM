package auth

// MailService specifies an API for sending emails via SMTP.
type MailService interface {

	// SendEmailAsHTML sends an HTML email using given template and set of template variables.
	SendEmailAsHTML(to string, subject string, templatePath string, templateData map[string]interface{}) error
}
