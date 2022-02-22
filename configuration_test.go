package truemail

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfiguration(t *testing.T) {
	validVerifierEmail, domain := pairRandomEmailDomain()

	t.Run("sets default configuration template", func(t *testing.T) {
		emptyStringSlice, emptyStringMap := []string(nil), map[string]string(nil)
		emailRegex, _ := newRegex(regexEmailPattern)
		smtpErrorBodyRegex, _ := newRegex(regexSMTPErrorBodyPattern)
		configurationAttr := ConfigurationAttr{VerifierEmail: validVerifierEmail}
		configuration, err := NewConfiguration(configurationAttr)

		assert.NoError(t, err)
		assert.Nil(t, configuration.ctx)
		assert.Equal(t, configurationAttr.VerifierEmail, configuration.VerifierEmail)
		assert.Equal(t, domain, configuration.VerifierDomain)
		assert.Equal(t, "smtp", configuration.ValidationTypeDefault)
		assert.Equal(t, defaultConnectionTimeout, configuration.ConnectionTimeout)
		assert.Equal(t, defaultResponseTimeout, configuration.ResponseTimeout)
		assert.Equal(t, defaultConnectionAttempts, configuration.ConnectionAttempts)
		assert.Equal(t, emptyStringSlice, configuration.WhitelistedDomains)
		assert.Equal(t, emptyStringSlice, configuration.BlacklistedDomains)
		assert.Equal(t, emptyStringSlice, configuration.BlacklistedMxIpAddresses)
		assert.Equal(t, emptyString, configuration.Dns)
		assert.Equal(t, emptyStringMap, configuration.ValidationTypeByDomain)
		assert.Equal(t, false, configuration.WhitelistValidation)
		assert.Equal(t, false, configuration.NotRfcMxLookupFlow)
		assert.Equal(t, defaultSmtpPort, configuration.SmtpPort)
		assert.Equal(t, false, configuration.SmtpFailFast)
		assert.Equal(t, false, configuration.SmtpSafeCheck)
		assert.Equal(t, emailRegex, configuration.EmailPattern)
		assert.Equal(t, smtpErrorBodyRegex, configuration.SmtpErrorBodyPattern)
	})

	t.Run("sets custom configuration template, custom DNS with port number", func(t *testing.T) {
		configurationAttr := ConfigurationAttr{
			ctx:                      context.TODO(),
			VerifierEmail:            validVerifierEmail,
			VerifierDomain:           randomDomain(),
			ValidationTypeDefault:    "mx",
			EmailPattern:             `\A.+@.+\z`,
			SmtpErrorBodyPattern:     `550{1}`,
			ConnectionTimeout:        randomPositiveNumber(),
			ResponseTimeout:          randomPositiveNumber(),
			ConnectionAttempts:       randomPositiveNumber(),
			WhitelistedDomains:       []string{randomDomain(), randomDomain()},
			BlacklistedDomains:       []string{randomDomain(), randomDomain()},
			BlacklistedMxIpAddresses: []string{randomIpAddress(), randomIpAddress()},
			Dns:                      randomDnsServer(),
			ValidationTypeByDomain:   map[string]string{randomDomain(): "regex"},
			WhitelistValidation:      true,
			NotRfcMxLookupFlow:       true,
			SmtpPort:                 randomPortNumber(),
			SmtpFailFast:             true,
			SmtpSafeCheck:            true,
		}
		emailRegex, _ := newRegex(configurationAttr.EmailPattern)
		smtpErrorBodyRegex, _ := newRegex(configurationAttr.SmtpErrorBodyPattern)
		configuration, err := NewConfiguration(configurationAttr)

		assert.NoError(t, err)
		assert.Equal(t, configurationAttr.ctx, configuration.ctx)
		assert.Equal(t, configurationAttr.VerifierEmail, configuration.VerifierEmail)
		assert.Equal(t, configurationAttr.VerifierDomain, configuration.VerifierDomain)
		assert.Equal(t, configurationAttr.ValidationTypeDefault, configuration.ValidationTypeDefault)
		assert.Equal(t, configurationAttr.ConnectionTimeout, configuration.ConnectionTimeout)
		assert.Equal(t, configurationAttr.ResponseTimeout, configuration.ResponseTimeout)
		assert.Equal(t, configurationAttr.WhitelistedDomains, configuration.WhitelistedDomains)
		assert.Equal(t, configurationAttr.BlacklistedDomains, configuration.BlacklistedDomains)
		assert.Equal(t, configurationAttr.BlacklistedMxIpAddresses, configuration.BlacklistedMxIpAddresses)
		assert.Equal(t, configurationAttr.Dns, configuration.Dns)
		assert.Equal(t, configurationAttr.ValidationTypeByDomain, configuration.ValidationTypeByDomain)
		assert.Equal(t, configurationAttr.WhitelistValidation, configuration.WhitelistValidation)
		assert.Equal(t, configurationAttr.NotRfcMxLookupFlow, configuration.NotRfcMxLookupFlow)
		assert.Equal(t, configurationAttr.SmtpPort, configuration.SmtpPort)
		assert.Equal(t, configurationAttr.SmtpFailFast, configuration.SmtpFailFast)
		assert.Equal(t, configurationAttr.SmtpSafeCheck, configuration.SmtpSafeCheck)
		assert.Equal(t, emailRegex, configuration.EmailPattern)
		assert.Equal(t, smtpErrorBodyRegex, configuration.SmtpErrorBodyPattern)
	})

	t.Run("sets custom configuration template, custom DNS without port number", func(t *testing.T) {
		configurationAttr := ConfigurationAttr{
			ctx:                      context.TODO(),
			VerifierEmail:            validVerifierEmail,
			VerifierDomain:           randomDomain(),
			ValidationTypeDefault:    randomValidationType(),
			EmailPattern:             `\A.+@.+\z`,
			SmtpErrorBodyPattern:     `550{1}`,
			ConnectionTimeout:        randomPositiveNumber(),
			ResponseTimeout:          randomPositiveNumber(),
			ConnectionAttempts:       randomPositiveNumber(),
			WhitelistedDomains:       []string{randomDomain(), randomDomain()},
			BlacklistedDomains:       []string{randomDomain(), randomDomain()},
			BlacklistedMxIpAddresses: []string{randomIpAddress(), randomIpAddress()},
			Dns:                      randomIpAddress(),
			ValidationTypeByDomain:   map[string]string{randomDomain(): randomValidationType()},
			WhitelistValidation:      true,
			NotRfcMxLookupFlow:       true,
			SmtpPort:                 randomPortNumber(),
			SmtpFailFast:             true,
			SmtpSafeCheck:            true,
		}
		emailRegex, _ := newRegex(configurationAttr.EmailPattern)
		smtpErrorBodyRegex, _ := newRegex(configurationAttr.SmtpErrorBodyPattern)
		configuration, err := NewConfiguration(configurationAttr)

		assert.NoError(t, err)
		assert.Equal(t, configurationAttr.ctx, configuration.ctx)
		assert.Equal(t, configurationAttr.VerifierEmail, configuration.VerifierEmail)
		assert.Equal(t, configurationAttr.VerifierDomain, configuration.VerifierDomain)
		assert.Equal(t, configurationAttr.ValidationTypeDefault, configuration.ValidationTypeDefault)
		assert.Equal(t, configurationAttr.ConnectionTimeout, configuration.ConnectionTimeout)
		assert.Equal(t, configurationAttr.ResponseTimeout, configuration.ResponseTimeout)
		assert.Equal(t, configurationAttr.ConnectionAttempts, configuration.ConnectionAttempts)
		assert.Equal(t, configurationAttr.WhitelistedDomains, configuration.WhitelistedDomains)
		assert.Equal(t, configurationAttr.BlacklistedDomains, configuration.BlacklistedDomains)
		assert.Equal(t, configurationAttr.BlacklistedMxIpAddresses, configuration.BlacklistedMxIpAddresses)
		assert.Equal(t, serverWithPortNumber(configurationAttr.Dns, defaultDnsPort), configuration.Dns)
		assert.Equal(t, configurationAttr.ValidationTypeByDomain, configuration.ValidationTypeByDomain)
		assert.Equal(t, configurationAttr.WhitelistValidation, configuration.WhitelistValidation)
		assert.Equal(t, configurationAttr.NotRfcMxLookupFlow, configuration.NotRfcMxLookupFlow)
		assert.Equal(t, configurationAttr.SmtpPort, configuration.SmtpPort)
		assert.Equal(t, configurationAttr.SmtpFailFast, configuration.SmtpFailFast)
		assert.Equal(t, configurationAttr.SmtpSafeCheck, configuration.SmtpSafeCheck)
		assert.Equal(t, emailRegex, configuration.EmailPattern)
		assert.Equal(t, smtpErrorBodyRegex, configuration.SmtpErrorBodyPattern)
	})

	t.Run("invalid verifier email", func(t *testing.T) {
		configurationAttr := ConfigurationAttr{VerifierEmail: "email@domain"}
		configuration, err := NewConfiguration(configurationAttr)
		errorMessage := fmt.Sprintf("%v is invalid verifier email", configurationAttr.VerifierEmail)

		assert.Nil(t, configuration)
		assert.EqualError(t, err, errorMessage)
	})

	t.Run("invalid verifier domain", func(t *testing.T) {
		configurationAttr := ConfigurationAttr{VerifierEmail: validVerifierEmail, VerifierDomain: "invalid_domain"}
		configuration, err := NewConfiguration(configurationAttr)
		errorMessage := fmt.Sprintf("%v is invalid verifier domain", configurationAttr.VerifierDomain)

		assert.Nil(t, configuration)
		assert.EqualError(t, err, errorMessage)
	})

	t.Run("invalid default validation type", func(t *testing.T) {
		configurationAttr := ConfigurationAttr{VerifierEmail: validVerifierEmail, ValidationTypeDefault: "invalid validation type"}
		configuration, err := NewConfiguration(configurationAttr)
		errorMessage := fmt.Sprintf(
			"%v is invalid default validation type, use one of these: [regex mx mx_blacklist smtp]",
			configurationAttr.ValidationTypeDefault,
		)

		assert.Nil(t, configuration)
		assert.EqualError(t, err, errorMessage)
	})

	t.Run("invalid connection timeout", func(t *testing.T) {
		configurationAttr := ConfigurationAttr{VerifierEmail: validVerifierEmail, ConnectionTimeout: -42}
		configuration, err := NewConfiguration(configurationAttr)
		errorMessage := fmt.Sprintf("%v should be a positive integer", configurationAttr.ConnectionTimeout)

		assert.Nil(t, configuration)
		assert.EqualError(t, err, errorMessage)
	})

	t.Run("invalid response timeout", func(t *testing.T) {
		configurationAttr := ConfigurationAttr{VerifierEmail: validVerifierEmail, ResponseTimeout: -42}
		configuration, err := NewConfiguration(configurationAttr)
		errorMessage := fmt.Sprintf("%v should be a positive integer", configurationAttr.ResponseTimeout)

		assert.Nil(t, configuration)
		assert.EqualError(t, err, errorMessage)
	})

	t.Run("invalid connection attempts", func(t *testing.T) {
		configurationAttr := ConfigurationAttr{VerifierEmail: validVerifierEmail, ConnectionAttempts: -42}
		configuration, err := NewConfiguration(configurationAttr)
		errorMessage := fmt.Sprintf("%v should be a positive integer", configurationAttr.ConnectionAttempts)

		assert.Nil(t, configuration)
		assert.EqualError(t, err, errorMessage)
	})

	t.Run("invalid SMTP port number", func(t *testing.T) {
		configurationAttr := ConfigurationAttr{VerifierEmail: validVerifierEmail, SmtpPort: -42}
		configuration, err := NewConfiguration(configurationAttr)
		errorMessage := fmt.Sprintf("%v should be a positive integer", configurationAttr.SmtpPort)

		assert.Nil(t, configuration)
		assert.EqualError(t, err, errorMessage)
	})

	t.Run("invalid whitelisted domains", func(t *testing.T) {
		configurationAttr := ConfigurationAttr{VerifierEmail: validVerifierEmail, WhitelistedDomains: []string{randomDomain(), "a"}}
		configuration, err := NewConfiguration(configurationAttr)
		errorMessage := fmt.Sprintf("%v is invalid domain name", configurationAttr.WhitelistedDomains[1])

		assert.Nil(t, configuration)
		assert.EqualError(t, err, errorMessage)
	})

	t.Run("invalid blacklisted domains", func(t *testing.T) {
		configurationAttr := ConfigurationAttr{VerifierEmail: validVerifierEmail, BlacklistedDomains: []string{randomDomain(), "b"}}
		configuration, err := NewConfiguration(configurationAttr)
		errorMessage := fmt.Sprintf("%v is invalid domain name", configurationAttr.BlacklistedDomains[1])

		assert.Nil(t, configuration)
		assert.EqualError(t, err, errorMessage)
	})

	t.Run("invalid blacklisted mx ip address", func(t *testing.T) {
		configurationAttr := ConfigurationAttr{VerifierEmail: validVerifierEmail, BlacklistedMxIpAddresses: []string{randomIpAddress(), "1.1.1.256:65536"}}
		configuration, err := NewConfiguration(configurationAttr)
		errorMessage := fmt.Sprintf("%v is invalid ip address", configurationAttr.BlacklistedMxIpAddresses[1])

		assert.Nil(t, configuration)
		assert.EqualError(t, err, errorMessage)
	})

	t.Run("invalid dns, wrong ip address", func(t *testing.T) {
		configurationAttr := ConfigurationAttr{VerifierEmail: validVerifierEmail, Dns: "1.1.1.256:65535"}
		configuration, err := NewConfiguration(configurationAttr)
		errorMessage := fmt.Sprintf("%v is invalid dns server", configurationAttr.Dns)

		assert.Nil(t, configuration)
		assert.EqualError(t, err, errorMessage)
	})

	t.Run("invalid dns, wrong port number", func(t *testing.T) {
		configurationAttr := ConfigurationAttr{VerifierEmail: validVerifierEmail, Dns: "2.2.2.2:65536"}
		configuration, err := NewConfiguration(configurationAttr)
		errorMessage := fmt.Sprintf("%v is invalid dns server", configurationAttr.Dns)

		assert.Nil(t, configuration)
		assert.EqualError(t, err, errorMessage)
	})

	t.Run("invalid dns, wrong ip address and port number", func(t *testing.T) {
		configurationAttr := ConfigurationAttr{VerifierEmail: validVerifierEmail, Dns: "1.1.1.256:65536"}
		configuration, err := NewConfiguration(configurationAttr)
		errorMessage := fmt.Sprintf("%v is invalid dns server", configurationAttr.Dns)

		assert.Nil(t, configuration)
		assert.EqualError(t, err, errorMessage)
	})

	t.Run("invalid validation type by domain, wrong domain", func(t *testing.T) {
		invalidDomain := "inavlid domain"
		configurationAttr := ConfigurationAttr{VerifierEmail: validVerifierEmail, ValidationTypeByDomain: map[string]string{randomDomain(): "regex", invalidDomain: "wrong_type"}}
		configuration, err := NewConfiguration(configurationAttr)
		errorMessage := fmt.Sprintf("%v is invalid domain name", invalidDomain)

		assert.Nil(t, configuration)
		assert.EqualError(t, err, errorMessage)
	})

	t.Run("invalid validation type by domain, wrong validation type", func(t *testing.T) {
		invalidType := "inavlid validation type"
		configurationAttr := ConfigurationAttr{VerifierEmail: validVerifierEmail, ValidationTypeByDomain: map[string]string{randomDomain(): "regex", randomDomain(): invalidType}}
		configuration, err := NewConfiguration(configurationAttr)
		errorMessage := fmt.Sprintf("%v is invalid default validation type, use one of these: [regex mx mx_blacklist smtp]", invalidType)

		assert.Nil(t, configuration)
		assert.EqualError(t, err, errorMessage)
	})

	t.Run("invalid email pattern", func(t *testing.T) {
		configurationAttr := ConfigurationAttr{VerifierEmail: validVerifierEmail, EmailPattern: `\K`}
		configuration, err := NewConfiguration(configurationAttr)
		errorMessage := fmt.Sprintf("error parsing regexp: invalid escape sequence: `%v`", configurationAttr.EmailPattern)

		assert.Nil(t, configuration)
		assert.EqualError(t, err, errorMessage)
	})

	t.Run("invalid smtp error body pattern", func(t *testing.T) {
		configurationAttr := ConfigurationAttr{VerifierEmail: validVerifierEmail, SmtpErrorBodyPattern: `\K`}
		configuration, err := NewConfiguration(configurationAttr)
		errorMessage := fmt.Sprintf("error parsing regexp: invalid escape sequence: `%v`", configurationAttr.SmtpErrorBodyPattern)

		assert.Nil(t, configuration)
		assert.EqualError(t, err, errorMessage)
	})
}
