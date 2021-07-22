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
		emailRegex, _ := newRegex(RegexEmailPattern)
		smtpErrorBodyRegex, _ := newRegex(RegexSMTPErrorBodyPattern)
		configurationAttr := ConfigurationAttr{verifierEmail: validVerifierEmail}
		configuration, err := NewConfiguration(configurationAttr)

		assert.NoError(t, err)
		assert.Nil(t, configuration.ctx)
		assert.Equal(t, configurationAttr.verifierEmail, configuration.VerifierEmail)
		assert.Equal(t, domain, configuration.VerifierDomain)
		assert.Equal(t, "smtp", configuration.ValidationTypeDefault)
		assert.Equal(t, DefaultConnectionTimeout, configuration.ConnectionTimeout)
		assert.Equal(t, DefaultResponseTimeout, configuration.ResponseTimeout)
		assert.Equal(t, DefaultConnectionAttempts, configuration.ConnectionAttempts)
		assert.Equal(t, emptyStringSlice, configuration.WhitelistedDomains)
		assert.Equal(t, emptyStringSlice, configuration.BlacklistedDomains)
		assert.Equal(t, emptyStringSlice, configuration.BlacklistedMxIpAddresses)
		assert.Equal(t, EmptyString, configuration.DNS)
		assert.Equal(t, emptyStringMap, configuration.ValidationTypeByDomain)
		assert.Equal(t, false, configuration.WhitelistValidation)
		assert.Equal(t, false, configuration.NotRfcMxLookupFlow)
		assert.Equal(t, false, configuration.SMTPFailFast)
		assert.Equal(t, false, configuration.SMTPSafeCheck)
		assert.Equal(t, emailRegex, configuration.EmailPattern)
		assert.Equal(t, smtpErrorBodyRegex, configuration.SMTPErrorBodyPattern)
	})

	t.Run("sets custom configuration template, custom DNS with port number", func(t *testing.T) {
		configurationAttr := ConfigurationAttr{
			ctx:                      context.TODO(),
			verifierEmail:            validVerifierEmail,
			verifierDomain:           randomDomain(),
			validationTypeDefault:    "mx",
			emailPattern:             `\A.+@.+\z`,
			smtpErrorBodyPattern:     `550{1}`,
			connectionTimeout:        3,
			responseTimeout:          4,
			connectionAttempts:       5,
			whitelistedDomains:       []string{randomDomain(), randomDomain()},
			blacklistedDomains:       []string{randomDomain(), randomDomain()},
			blacklistedMxIpAddresses: []string{randomIpAddress(), randomIpAddress()},
			dns:                      randomDnsServer(),
			validationTypeByDomain:   map[string]string{randomDomain(): "regex"},
			whitelistValidation:      true,
			notRfcMxLookupFlow:       true,
			smtpFailFast:             true,
			smtpSafeCheck:            true,
		}
		emailRegex, _ := newRegex(configurationAttr.emailPattern)
		smtpErrorBodyRegex, _ := newRegex(configurationAttr.smtpErrorBodyPattern)
		configuration, err := NewConfiguration(configurationAttr)

		assert.NoError(t, err)
		assert.Equal(t, configurationAttr.ctx, configuration.ctx)
		assert.Equal(t, configurationAttr.verifierEmail, configuration.VerifierEmail)
		assert.Equal(t, configurationAttr.verifierDomain, configuration.VerifierDomain)
		assert.Equal(t, configurationAttr.validationTypeDefault, configuration.ValidationTypeDefault)
		assert.Equal(t, configurationAttr.connectionTimeout, configuration.ConnectionTimeout)
		assert.Equal(t, configurationAttr.responseTimeout, configuration.ResponseTimeout)
		assert.Equal(t, configurationAttr.connectionAttempts, configuration.ConnectionAttempts)
		assert.Equal(t, configurationAttr.whitelistedDomains, configuration.WhitelistedDomains)
		assert.Equal(t, configurationAttr.blacklistedDomains, configuration.BlacklistedDomains)
		assert.Equal(t, configurationAttr.blacklistedMxIpAddresses, configuration.BlacklistedMxIpAddresses)
		assert.Equal(t, configurationAttr.dns, configuration.DNS)
		assert.Equal(t, configurationAttr.validationTypeByDomain, configuration.ValidationTypeByDomain)
		assert.Equal(t, configurationAttr.whitelistValidation, configuration.WhitelistValidation)
		assert.Equal(t, configurationAttr.notRfcMxLookupFlow, configuration.NotRfcMxLookupFlow)
		assert.Equal(t, configurationAttr.smtpFailFast, configuration.SMTPFailFast)
		assert.Equal(t, configurationAttr.smtpSafeCheck, configuration.SMTPSafeCheck)
		assert.Equal(t, emailRegex, configuration.EmailPattern)
		assert.Equal(t, smtpErrorBodyRegex, configuration.SMTPErrorBodyPattern)
	})

	t.Run("sets custom configuration template, custom DNS without port number", func(t *testing.T) {
		configurationAttr := ConfigurationAttr{
			ctx:                      context.TODO(),
			verifierEmail:            validVerifierEmail,
			verifierDomain:           randomDomain(),
			validationTypeDefault:    randomValidationType(),
			emailPattern:             `\A.+@.+\z`,
			smtpErrorBodyPattern:     `550{1}`,
			connectionTimeout:        3,
			responseTimeout:          4,
			connectionAttempts:       5,
			whitelistedDomains:       []string{randomDomain(), randomDomain()},
			blacklistedDomains:       []string{randomDomain(), randomDomain()},
			blacklistedMxIpAddresses: []string{randomIpAddress(), randomIpAddress()},
			dns:                      randomIpAddress(),
			validationTypeByDomain:   map[string]string{randomDomain(): randomValidationType()},
			whitelistValidation:      true,
			notRfcMxLookupFlow:       true,
			smtpFailFast:             true,
			smtpSafeCheck:            true,
		}
		emailRegex, _ := newRegex(configurationAttr.emailPattern)
		smtpErrorBodyRegex, _ := newRegex(configurationAttr.smtpErrorBodyPattern)
		configuration, err := NewConfiguration(configurationAttr)

		assert.NoError(t, err)
		assert.Equal(t, configurationAttr.ctx, configuration.ctx)
		assert.Equal(t, configurationAttr.verifierEmail, configuration.VerifierEmail)
		assert.Equal(t, configurationAttr.verifierDomain, configuration.VerifierDomain)
		assert.Equal(t, configurationAttr.validationTypeDefault, configuration.ValidationTypeDefault)
		assert.Equal(t, configurationAttr.connectionTimeout, configuration.ConnectionTimeout)
		assert.Equal(t, configurationAttr.responseTimeout, configuration.ResponseTimeout)
		assert.Equal(t, configurationAttr.connectionAttempts, configuration.ConnectionAttempts)
		assert.Equal(t, configurationAttr.whitelistedDomains, configuration.WhitelistedDomains)
		assert.Equal(t, configurationAttr.blacklistedDomains, configuration.BlacklistedDomains)
		assert.Equal(t, configurationAttr.blacklistedMxIpAddresses, configuration.BlacklistedMxIpAddresses)
		assert.Equal(t, configurationAttr.dns+":"+DefaultDnsPort, configuration.DNS)
		assert.Equal(t, configurationAttr.validationTypeByDomain, configuration.ValidationTypeByDomain)
		assert.Equal(t, configurationAttr.whitelistValidation, configuration.WhitelistValidation)
		assert.Equal(t, configurationAttr.notRfcMxLookupFlow, configuration.NotRfcMxLookupFlow)
		assert.Equal(t, configurationAttr.smtpFailFast, configuration.SMTPFailFast)
		assert.Equal(t, configurationAttr.smtpSafeCheck, configuration.SMTPSafeCheck)
		assert.Equal(t, emailRegex, configuration.EmailPattern)
		assert.Equal(t, smtpErrorBodyRegex, configuration.SMTPErrorBodyPattern)
	})

	t.Run("invalid verifier email", func(t *testing.T) {
		configurationAttr := ConfigurationAttr{verifierEmail: "email@domain"}
		configuration, err := NewConfiguration(configurationAttr)
		errorMessage := fmt.Sprintf("%v is invalid verifier email", configurationAttr.verifierEmail)

		assert.Nil(t, configuration)
		assert.EqualError(t, err, errorMessage)
	})

	t.Run("invalid verifier domain", func(t *testing.T) {
		configurationAttr := ConfigurationAttr{verifierEmail: validVerifierEmail, verifierDomain: "invalid_domain"}
		configuration, err := NewConfiguration(configurationAttr)
		errorMessage := fmt.Sprintf("%v is invalid verifier domain", configurationAttr.verifierDomain)

		assert.Nil(t, configuration)
		assert.EqualError(t, err, errorMessage)
	})

	t.Run("invalid default validation type", func(t *testing.T) {
		configurationAttr := ConfigurationAttr{verifierEmail: validVerifierEmail, validationTypeDefault: "invalid validation type"}
		configuration, err := NewConfiguration(configurationAttr)
		errorMessage := fmt.Sprintf("%v is invalid default validation type, use one of these: [regex mx mx_blacklist smtp]", configurationAttr.validationTypeDefault)

		assert.Nil(t, configuration)
		assert.EqualError(t, err, errorMessage)
	})

	t.Run("invalid connection timeout", func(t *testing.T) {
		configurationAttr := ConfigurationAttr{verifierEmail: validVerifierEmail, connectionTimeout: -42}
		configuration, err := NewConfiguration(configurationAttr)
		errorMessage := fmt.Sprintf("%v should be a positive integer", configurationAttr.connectionTimeout)

		assert.Nil(t, configuration)
		assert.EqualError(t, err, errorMessage)
	})

	t.Run("invalid response timeout", func(t *testing.T) {
		configurationAttr := ConfigurationAttr{verifierEmail: validVerifierEmail, responseTimeout: -42}
		configuration, err := NewConfiguration(configurationAttr)
		errorMessage := fmt.Sprintf("%v should be a positive integer", configurationAttr.responseTimeout)

		assert.Nil(t, configuration)
		assert.EqualError(t, err, errorMessage)
	})

	t.Run("invalid connection attempts", func(t *testing.T) {
		configurationAttr := ConfigurationAttr{verifierEmail: validVerifierEmail, connectionAttempts: -42}
		configuration, err := NewConfiguration(configurationAttr)
		errorMessage := fmt.Sprintf("%v should be a positive integer", configurationAttr.connectionAttempts)

		assert.Nil(t, configuration)
		assert.EqualError(t, err, errorMessage)
	})

	t.Run("invalid whitelisted domains", func(t *testing.T) {
		configurationAttr := ConfigurationAttr{verifierEmail: validVerifierEmail, whitelistedDomains: []string{randomDomain(), "a"}}
		configuration, err := NewConfiguration(configurationAttr)
		errorMessage := fmt.Sprintf("%v is invalid domain name", configurationAttr.whitelistedDomains[1])

		assert.Nil(t, configuration)
		assert.EqualError(t, err, errorMessage)
	})

	t.Run("invalid blacklisted domains", func(t *testing.T) {
		configurationAttr := ConfigurationAttr{verifierEmail: validVerifierEmail, blacklistedDomains: []string{randomDomain(), "b"}}
		configuration, err := NewConfiguration(configurationAttr)
		errorMessage := fmt.Sprintf("%v is invalid domain name", configurationAttr.blacklistedDomains[1])

		assert.Nil(t, configuration)
		assert.EqualError(t, err, errorMessage)
	})

	t.Run("invalid blacklisted mx ip address", func(t *testing.T) {
		configurationAttr := ConfigurationAttr{verifierEmail: validVerifierEmail, blacklistedMxIpAddresses: []string{randomIpAddress(), "1.1.1.256:65536"}}
		configuration, err := NewConfiguration(configurationAttr)
		errorMessage := fmt.Sprintf("%v is invalid ip address", configurationAttr.blacklistedMxIpAddresses[1])

		assert.Nil(t, configuration)
		assert.EqualError(t, err, errorMessage)
	})

	t.Run("invalid dns, wrong ip address", func(t *testing.T) {
		configurationAttr := ConfigurationAttr{verifierEmail: validVerifierEmail, dns: "1.1.1.256:65535"}
		configuration, err := NewConfiguration(configurationAttr)
		errorMessage := fmt.Sprintf("%v is invalid dns server", configurationAttr.dns)

		assert.Nil(t, configuration)
		assert.EqualError(t, err, errorMessage)
	})

	t.Run("invalid dns, wrong port number", func(t *testing.T) {
		configurationAttr := ConfigurationAttr{verifierEmail: validVerifierEmail, dns: "2.2.2.2:65536"}
		configuration, err := NewConfiguration(configurationAttr)
		errorMessage := fmt.Sprintf("%v is invalid dns server", configurationAttr.dns)

		assert.Nil(t, configuration)
		assert.EqualError(t, err, errorMessage)
	})

	t.Run("invalid dns, wrong ip address and port number", func(t *testing.T) {
		configurationAttr := ConfigurationAttr{verifierEmail: validVerifierEmail, dns: "1.1.1.256:65536"}
		configuration, err := NewConfiguration(configurationAttr)
		errorMessage := fmt.Sprintf("%v is invalid dns server", configurationAttr.dns)

		assert.Nil(t, configuration)
		assert.EqualError(t, err, errorMessage)
	})

	t.Run("invalid validation type by domain, wrong domain", func(t *testing.T) {
		invalidDomain := "inavlid domain"
		configurationAttr := ConfigurationAttr{verifierEmail: validVerifierEmail, validationTypeByDomain: map[string]string{randomDomain(): "regex", invalidDomain: "wrong_type"}}
		configuration, err := NewConfiguration(configurationAttr)
		errorMessage := fmt.Sprintf("%v is invalid domain name", invalidDomain)

		assert.Nil(t, configuration)
		assert.EqualError(t, err, errorMessage)
	})

	t.Run("invalid validation type by domain, wrong validation type", func(t *testing.T) {
		invalidType := "inavlid validation type"
		configurationAttr := ConfigurationAttr{verifierEmail: validVerifierEmail, validationTypeByDomain: map[string]string{randomDomain(): "regex", randomDomain(): invalidType}}
		configuration, err := NewConfiguration(configurationAttr)
		errorMessage := fmt.Sprintf("%v is invalid default validation type, use one of these: [regex mx mx_blacklist smtp]", invalidType)

		assert.Nil(t, configuration)
		assert.EqualError(t, err, errorMessage)
	})

	t.Run("invalid email pattern", func(t *testing.T) {
		configurationAttr := ConfigurationAttr{verifierEmail: validVerifierEmail, emailPattern: `\K`}
		configuration, err := NewConfiguration(configurationAttr)
		errorMessage := fmt.Sprintf("error parsing regexp: invalid escape sequence: `%v`", configurationAttr.emailPattern)

		assert.Nil(t, configuration)
		assert.EqualError(t, err, errorMessage)
	})

	t.Run("invalid smtp error body pattern", func(t *testing.T) {
		configurationAttr := ConfigurationAttr{verifierEmail: validVerifierEmail, smtpErrorBodyPattern: `\K`}
		configuration, err := NewConfiguration(configurationAttr)
		errorMessage := fmt.Sprintf("error parsing regexp: invalid escape sequence: `%v`", configurationAttr.smtpErrorBodyPattern)

		assert.Nil(t, configuration)
		assert.EqualError(t, err, errorMessage)
	})
}
