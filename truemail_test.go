package truemail

import (
	"fmt"
	"net"
	"testing"

	"github.com/foxcpp/go-mockdns"
	smtpmock "github.com/mocktools/go-smtp-mock"
	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	email, domain := pairRandomEmailDomain()
	resolvedHostNameByMxReord := randomDnsHostName()
	blacklistedEmail, blacklistedDomain := pairRandomEmailDomain()
	blacklistedResolvedHostNameByMxReord, blacklistedMxIpAddress, blacklistedVerifierEmail := randomDnsHostName(), randomIpAddress(), randomEmail()

	server := startSmtpMock(smtpmock.ConfigurationAttr{BlacklistedMailfromEmails: []string{blacklistedVerifierEmail}})
	portNumber := server.PortNumber
	defer func() { _ = server.Stop() }()

	dns := runMockDnsServer(
		map[string]mockdns.Zone{
			toDnsHostName(punycodeDomain(domain)): {
				MX: []net.MX{
					{Host: resolvedHostNameByMxReord, Pref: uint16(5)},
				},
			},
			resolvedHostNameByMxReord: {
				A: []string{localhostIPv4Address},
			},
			toDnsHostName(punycodeDomain(blacklistedDomain)): {
				MX: []net.MX{
					{Host: blacklistedResolvedHostNameByMxReord, Pref: uint16(5)},
				},
			},
			blacklistedResolvedHostNameByMxReord: {
				A: []string{blacklistedMxIpAddress},
			},
		},
	)

	for _, validValidationType := range availableValidationTypes() {
		t.Run(validValidationType+" valid validation type", func(t *testing.T) {
			configuration, _ := NewConfiguration(
				ConfigurationAttr{
					VerifierEmail: randomEmail(),
					Dns:           dns,
					SmtpPort:      portNumber,
				},
			)
			validatorResult, err := Validate(email, configuration, validValidationType)

			assert.NoError(t, err)
			assert.Equal(t, email, validatorResult.Email)
			assert.NotSame(t, configuration, validatorResult.Configuration)
			assert.Equal(t, validValidationType, validatorResult.ValidationType)
			assert.Equal(t, usedValidationsByType(validValidationType), validatorResult.usedValidations)
			assert.True(t, validatorResult.Success)
		})
	}

	t.Run("succesful validation, default validation type specified in configuration", func(t *testing.T) {
		email, specifiedValidationTypeByDefault := randomEmail(), validationTypeRegex
		configuration, _ := NewConfiguration(
			ConfigurationAttr{
				VerifierEmail:         randomEmail(),
				ValidationTypeDefault: specifiedValidationTypeByDefault,
			},
		)
		validatorResult, _ := Validate(email, configuration)

		assert.True(t, validatorResult.Success)
		assert.True(t, validatorResult.isPassFromDomainListMatch)
		assert.Equal(t, []string{specifiedValidationTypeByDefault}, validatorResult.usedValidations)
	})

	t.Run("invalid validation type", func(t *testing.T) {
		invalidValidationType := "invalid type"
		errorMessage := fmt.Sprintf("%s is invalid validation type, use one of these: [regex mx mx_blacklist smtp]", invalidValidationType)
		_, err := Validate(randomEmail(), createConfiguration(), invalidValidationType)
		assert.EqualError(t, err, errorMessage)
	})

	t.Run("Whitelist/Blacklist validation successful", func(t *testing.T) {
		email, domain := pairRandomEmailDomain()
		configuration, _ := NewConfiguration(
			ConfigurationAttr{
				VerifierEmail:      randomEmail(),
				WhitelistedDomains: []string{domain},
			},
		)
		validatorResult, _ := Validate(email, configuration)

		assert.True(t, validatorResult.Success)
		assert.False(t, validatorResult.isPassFromDomainListMatch)
		assert.Empty(t, validatorResult.usedValidations)
	})

	t.Run("Whitelist/Blacklist validation passes to next validation level", func(t *testing.T) {
		configuration, _ := NewConfiguration(
			ConfigurationAttr{
				VerifierEmail:       randomEmail(),
				WhitelistedDomains:  []string{domain},
				WhitelistValidation: true,
				Dns:                 dns,
				SmtpPort:            portNumber,
			},
		)
		validatorResult, _ := Validate(email, configuration)

		assert.True(t, validatorResult.Success)
		assert.True(t, validatorResult.isPassFromDomainListMatch)
		assert.Equal(t, usedValidationsByType(validationTypeSmtp), validatorResult.usedValidations)
	})

	t.Run("Whitelist/Blacklist validation fails", func(t *testing.T) {
		configuration, _ := NewConfiguration(
			ConfigurationAttr{
				VerifierEmail:      randomEmail(),
				BlacklistedDomains: []string{domain},
			},
		)
		validatorResult, _ := Validate(email, configuration)

		assert.False(t, validatorResult.Success)
		assert.False(t, validatorResult.isPassFromDomainListMatch)
		assert.Empty(t, validatorResult.usedValidations)
	})

	t.Run("Mx blacklist validation fails", func(t *testing.T) {
		configuration, _ := NewConfiguration(
			ConfigurationAttr{
				VerifierEmail:            randomEmail(),
				BlacklistedMxIpAddresses: []string{blacklistedMxIpAddress},
				Dns:                      dns,
			},
		)
		validatorResult, _ := Validate(blacklistedEmail, configuration)

		assert.False(t, validatorResult.Success)
		assert.Equal(t, usedValidationsByType(validationTypeMxBlacklist), validatorResult.usedValidations)
	})

	t.Run("SMTP validation fails", func(t *testing.T) {
		configuration, _ := NewConfiguration(
			ConfigurationAttr{
				VerifierEmail: blacklistedVerifierEmail,
				Dns:           dns,
				SmtpPort:      portNumber,
			},
		)
		validatorResult, _ := Validate(email, configuration)

		assert.False(t, validatorResult.Success)
		assert.Equal(t, usedValidationsByType(validationTypeSmtp), validatorResult.usedValidations)
	})
}

func TestIsValid(t *testing.T) {
	email, domain := pairRandomEmailDomain()
	resolvedHostNameByMxReord := randomDnsHostName()
	nonExistentEmail, nonExistentDomain := pairRandomEmailDomain()
	nonExistentResolvedHostNameByMxReord := randomDnsHostName()

	server := startSmtpMock(smtpmock.ConfigurationAttr{NotRegisteredEmails: []string{nonExistentEmail}})
	portNumber := server.PortNumber
	defer func() { _ = server.Stop() }()

	dns := runMockDnsServer(
		map[string]mockdns.Zone{
			toDnsHostName(punycodeDomain(domain)): {
				MX: []net.MX{
					{Host: resolvedHostNameByMxReord, Pref: uint16(5)},
				},
			},
			resolvedHostNameByMxReord: {
				A: []string{localhostIPv4Address},
			},
			toDnsHostName(punycodeDomain(nonExistentDomain)): {
				MX: []net.MX{
					{Host: nonExistentResolvedHostNameByMxReord, Pref: uint16(5)},
				},
			},
			nonExistentResolvedHostNameByMxReord: {
				A: []string{localhostIPv4Address},
			},
		},
	)

	configuration, _ := NewConfiguration(
		ConfigurationAttr{
			VerifierEmail: randomEmail(),
			Dns:           dns,
			SmtpPort:      portNumber,
		},
	)

	t.Run("when succesful validation, default validation type specified in configuration", func(t *testing.T) {
		assert.True(t, IsValid(email, configuration))
	})

	t.Run("when failed validation, default validation type specified in configuration", func(t *testing.T) {
		assert.False(t, IsValid(nonExistentEmail, configuration))
	})

	t.Run("when succesful validation, specified validation type", func(t *testing.T) {
		assert.True(t, IsValid(randomEmail(), createConfiguration(), validationTypeRegex))
	})

	t.Run("when failed validation", func(t *testing.T) {
		assert.False(t, IsValid("invalid@email", createConfiguration(), validationTypeRegex))
	})

	t.Run("when invalid validation type", func(t *testing.T) {
		assert.False(t, IsValid(randomEmail(), createConfiguration(), "invalidValidationType"))
	})
}
