package truemail

// SMTP validation builder entities interface
type builder interface {
	newSmtpRequest(int, string, string, *Configuration) *SmtpRequest
	newSmtpClient(*SmtpRequestConfiguration) client
}

// SMTP entities builder structure
type smtpBuilder struct{}

// interface implementation

// SMTP request builder. Returns pointer to configured new SMTP request structure
func (builder *smtpBuilder) newSmtpRequest(attempts int, targetEmail, targetHostAddress string, configuration *Configuration) *SmtpRequest {
	return &SmtpRequest{
		Attempts:      attempts,
		Email:         targetEmail,
		Host:          targetHostAddress,
		Configuration: newSmtpRequestConfiguration(configuration, targetEmail, targetHostAddress),
		Response:      new(SmtpResponse),
	}
}

// SMTP client builder. Returns pointer to configured new SMTP client
func (builder *smtpBuilder) newSmtpClient(configuration *SmtpRequestConfiguration) client {
	return newSmtpClient(configuration)
}
