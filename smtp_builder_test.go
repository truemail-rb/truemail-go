package truemail

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSmtpBuilderNewSmtpRequest(t *testing.T) {
	t.Run("creates new configured SMTP request", func(t *testing.T) {
		attempts, targetEmail, targetHostAddress, configuration := randomPositiveNumber(), randomEmail(), randomIpAddress(), createConfiguration()
		smtpRequestConfiguration := newSmtpRequestConfiguration(configuration, targetEmail, targetHostAddress)
		smtpRequest := new(smtpBuilder).newSmtpRequest(attempts, targetEmail, targetHostAddress, configuration)

		assert.Equal(t, attempts, smtpRequest.Attempts)
		assert.Equal(t, targetEmail, smtpRequest.Email)
		assert.Equal(t, targetHostAddress, smtpRequest.Host)
		assert.Equal(t, smtpRequestConfiguration, smtpRequest.Configuration)
		assert.Equal(t, new(SmtpResponse), smtpRequest.Response)
	})
}

func TestSmtpBuilderNewSmtpClient(t *testing.T) {
	t.Run("creates new configured SMTP client", func(t *testing.T) {
		smtpRequestConfiguration := newSmtpRequestConfiguration(createConfiguration(), randomEmail(), randomIpAddress())

		assert.Equal(t, newSmtpClient(smtpRequestConfiguration), new(smtpBuilder).newSmtpClient(smtpRequestConfiguration))
	})
}
