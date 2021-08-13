package truemail

import (
	"fmt"
	"net"
	"testing"

	"github.com/foxcpp/go-mockdns"
	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	// TODO: change to integration tests when validationSmtp.check() will be implemented
	for _, validValidationType := range availableValidationTypes() {
		t.Run(validValidationType+" valid validation type", func(t *testing.T) {
			email, domain := pairRandomEmailDomain()
			resolvedHostNameByMxReord := randomDnsHostName()
			dnsRecords := map[string]mockdns.Zone{
				toDnsHostName(punycodeDomain(domain)): {
					MX: []net.MX{
						{Host: resolvedHostNameByMxReord, Pref: uint16(5)},
					},
				},
				resolvedHostNameByMxReord: {
					A: []string{randomIpAddress()},
				},
			}
			configuration, _ := NewConfiguration(
				ConfigurationAttr{
					verifierEmail: randomEmail(),
					dns:           runMockDnsServer(dnsRecords),
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
		email, specifiedValidationTypeByDefault := randomEmail(), ValidationTypeRegex
		configuration, _ := NewConfiguration(
			ConfigurationAttr{
				verifierEmail:         randomEmail(),
				validationTypeDefault: specifiedValidationTypeByDefault,
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
				verifierEmail:      randomEmail(),
				whitelistedDomains: []string{domain},
			},
		)
		validatorResult, _ := Validate(email, configuration)

		assert.True(t, validatorResult.Success)
		assert.False(t, validatorResult.isPassFromDomainListMatch)
		assert.Empty(t, validatorResult.usedValidations)
	})

	t.Run("Whitelist/Blacklist validation passes to next validation level", func(t *testing.T) {
		email, domain := pairRandomEmailDomain()
		resolvedHostNameByMxReord := randomDnsHostName()
		dnsRecords := map[string]mockdns.Zone{
			toDnsHostName(punycodeDomain(domain)): {
				MX: []net.MX{
					{Host: resolvedHostNameByMxReord, Pref: uint16(5)},
				},
			},
			resolvedHostNameByMxReord: {
				A: []string{randomIpAddress()},
			},
		}
		configuration, _ := NewConfiguration(
			ConfigurationAttr{
				verifierEmail:       randomEmail(),
				whitelistedDomains:  []string{domain},
				whitelistValidation: true,
				dns:                 runMockDnsServer(dnsRecords),
			},
		)
		validatorResult, _ := Validate(email, configuration, ValidationTypeMx)

		assert.True(t, validatorResult.Success)
		assert.True(t, validatorResult.isPassFromDomainListMatch)
		assert.Equal(t, usedValidationsByType(ValidationTypeMx), validatorResult.usedValidations)
	})

	t.Run("Whitelist/Blacklist validation fails", func(t *testing.T) {
		email, domain := pairRandomEmailDomain()
		configuration, _ := NewConfiguration(
			ConfigurationAttr{
				verifierEmail:      randomEmail(),
				blacklistedDomains: []string{domain},
			},
		)
		validatorResult, _ := Validate(email, configuration)

		assert.False(t, validatorResult.Success)
		assert.False(t, validatorResult.isPassFromDomainListMatch)
		assert.Empty(t, validatorResult.usedValidations)
	})

	t.Run("Mx blacklist validation fails", func(t *testing.T) {
		email, domain := pairRandomEmailDomain()
		resolvedHostNameByMxReord, blacklistedMxIpAddress := randomDnsHostName(), randomIpAddress()
		dnsRecords := map[string]mockdns.Zone{
			toDnsHostName(punycodeDomain(domain)): {
				MX: []net.MX{
					{Host: resolvedHostNameByMxReord, Pref: uint16(5)},
				},
			},
			resolvedHostNameByMxReord: {
				A: []string{blacklistedMxIpAddress},
			},
		}
		configuration, _ := NewConfiguration(
			ConfigurationAttr{
				verifierEmail:            randomEmail(),
				blacklistedMxIpAddresses: []string{blacklistedMxIpAddress},
				dns:                      runMockDnsServer(dnsRecords),
			},
		)
		validatorResult, _ := Validate(email, configuration, ValidationTypeMxBlacklist)

		assert.False(t, validatorResult.Success)
		assert.Equal(t, usedValidationsByType(ValidationTypeMxBlacklist), validatorResult.usedValidations)
	})
}

func TestIsValid(t *testing.T) {
	// TODO: uncomment when validationSmtp.check() will be implemented
	// t.Run("when succesful validation, default validation type specified in configuration", func(t *testing.T) {
	// 	assert.True(t, IsValid(randomEmail(), createConfiguration()))
	// })

	t.Run("when succesful validation, specified validation type", func(t *testing.T) {
		assert.True(t, IsValid(randomEmail(), createConfiguration(), ValidationTypeRegex))
	})

	t.Run("when failure validation", func(t *testing.T) {
		assert.False(t, IsValid("invalid@email", createConfiguration(), ValidationTypeRegex))
	})

	t.Run("when invalid validation type", func(t *testing.T) {
		assert.False(t, IsValid(randomEmail(), createConfiguration(), "invalidValidationType"))
	})
}
