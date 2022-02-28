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

		assert.Equal(t, attempts, smtpRequest.attempts)
		assert.Equal(t, targetEmail, smtpRequest.email)
		assert.Equal(t, targetHostAddress, smtpRequest.host)
		assert.Equal(t, smtpRequestConfiguration, smtpRequest.configuration)
		assert.Equal(t, new(smtpResponse), smtpRequest.response)
	})
}

func TestSmtpBuilderNewSmtpClient(t *testing.T) {
	t.Run("creates new configured SMTP client", func(t *testing.T) {
		smtpRequestConfiguration := newSmtpRequestConfiguration(createConfiguration(), randomEmail(), randomIpAddress())

		assert.Equal(t, newSmtpClient(smtpRequestConfiguration), new(smtpBuilder).newSmtpClient(smtpRequestConfiguration))
	})
}
