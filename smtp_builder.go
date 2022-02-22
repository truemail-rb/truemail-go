package truemail

// SMTP validation builder entities interface
type builder interface {
	newSmtpRequest(int, string, string, *configuration) *smtpRequest
	newSmtpClient(*smtpRequestConfiguration) client
}

// SMTP entities builder structure
type smtpBuilder struct{}

// interface implementation

// SMTP request builder. Returns pointer to configured new SMTP request structure
func (builder *smtpBuilder) newSmtpRequest(attempts int, targetEmail, targetHostAddress string, configuration *configuration) *smtpRequest {
	return &smtpRequest{
		attempts:      attempts,
		email:         targetEmail,
		host:          targetHostAddress,
		configuration: newSmtpRequestConfiguration(configuration, targetEmail, targetHostAddress),
		response:      new(smtpResponse),
	}
}

// SMTP client builder. Returns pointer to configured new SMTP client
func (builder *smtpBuilder) newSmtpClient(configuration *smtpRequestConfiguration) client {
	return newSmtpClient(configuration)
}
