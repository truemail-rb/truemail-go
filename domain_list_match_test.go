package truemail

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidationDomainListMatchCheck(t *testing.T) {
	email, domain := pairRandomEmailDomain()

	t.Run("whitelist case, email is in whitelist", func(t *testing.T) {
		configuration, _ := NewConfiguration(
			ConfigurationAttr{
				VerifierEmail:      randomEmail(),
				WhitelistedDomains: []string{domain},
			},
		)
		validatorResult := runDomainListMatchValidation(email, configuration)

		assert.True(t, validatorResult.Success)
		assert.False(t, validatorResult.isPassFromDomainListMatch)
		assert.Equal(t, domainListMatchWhitelist, validatorResult.ValidationType)
		assert.Equal(t, domain, validatorResult.Domain)
	})

	t.Run("blacklist case, email is in blacklist", func(t *testing.T) {
		configuration, _ := NewConfiguration(
			ConfigurationAttr{
				VerifierEmail:      randomEmail(),
				BlacklistedDomains: []string{domain},
			},
		)
		validatorResult := runDomainListMatchValidation(email, configuration)

		assert.False(t, validatorResult.Success)
		assert.False(t, validatorResult.isPassFromDomainListMatch)
		assert.Equal(t, domainListMatchBlacklist, validatorResult.ValidationType)
		assert.Equal(t, domain, validatorResult.Domain)
	})

	t.Run("whitelist/blackist case, email is not in both lists", func(t *testing.T) {
		configuration, _ := NewConfiguration(
			ConfigurationAttr{
				VerifierEmail: randomEmail(),
			},
		)
		validationType := randomValidationType()
		validatorResult := runDomainListMatchValidation(email, configuration, validationType)

		assert.True(t, validatorResult.Success)
		assert.True(t, validatorResult.isPassFromDomainListMatch)
		assert.Equal(t, validationType, validatorResult.ValidationType)
		assert.Equal(t, domain, validatorResult.Domain)
	})

	t.Run("whitelist/blackist case, email is in both lists", func(t *testing.T) {
		configuration, _ := NewConfiguration(
			ConfigurationAttr{
				VerifierEmail:      randomEmail(),
				WhitelistedDomains: []string{domain},
				BlacklistedDomains: []string{domain},
			},
		)
		validatorResult := runDomainListMatchValidation(email, configuration)

		assert.False(t, validatorResult.Success)
		assert.False(t, validatorResult.isPassFromDomainListMatch)
		assert.Equal(t, domainListMatchBlacklist, validatorResult.ValidationType)
		assert.Equal(t, domain, validatorResult.Domain)
	})

	t.Run("whitelist validation case, email is in whitelist", func(t *testing.T) {
		configuration, _ := NewConfiguration(
			ConfigurationAttr{
				VerifierEmail:       randomEmail(),
				WhitelistValidation: true,
				WhitelistedDomains:  []string{domain},
			},
		)
		validationType := randomValidationType()
		validatorResult := runDomainListMatchValidation(email, configuration, validationType)

		assert.True(t, validatorResult.Success)
		assert.True(t, validatorResult.isPassFromDomainListMatch)
		assert.Equal(t, validationType, validatorResult.ValidationType)
		assert.Equal(t, domain, validatorResult.Domain)
	})

	t.Run("whitelist validation case, email is not in whitelist/blacklist", func(t *testing.T) {
		configuration, _ := NewConfiguration(
			ConfigurationAttr{
				VerifierEmail:       randomEmail(),
				WhitelistValidation: true,
			},
		)
		validatorResult := runDomainListMatchValidation(email, configuration)

		assert.False(t, validatorResult.Success)
		assert.False(t, validatorResult.isPassFromDomainListMatch)
		assert.Equal(t, domainListMatchBlacklist, validatorResult.ValidationType)
		assert.Equal(t, domain, validatorResult.Domain)
	})

	t.Run("whitelist validation case, email is in blacklist", func(t *testing.T) {
		configuration, _ := NewConfiguration(
			ConfigurationAttr{
				VerifierEmail:       randomEmail(),
				WhitelistValidation: true,
				BlacklistedDomains:  []string{domain},
			},
		)
		validatorResult := runDomainListMatchValidation(email, configuration)

		assert.False(t, validatorResult.Success)
		assert.False(t, validatorResult.isPassFromDomainListMatch)
		assert.Equal(t, domainListMatchBlacklist, validatorResult.ValidationType)
		assert.Equal(t, domain, validatorResult.Domain)
	})

	t.Run("whitelist validation case, email is in both lists", func(t *testing.T) {
		configuration, _ := NewConfiguration(
			ConfigurationAttr{
				VerifierEmail:       randomEmail(),
				WhitelistValidation: true,
				WhitelistedDomains:  []string{domain},
				BlacklistedDomains:  []string{domain},
			},
		)
		validatorResult := runDomainListMatchValidation(email, configuration)

		assert.False(t, validatorResult.Success)
		assert.False(t, validatorResult.isPassFromDomainListMatch)
		assert.Equal(t, domainListMatchBlacklist, validatorResult.ValidationType)
		assert.Equal(t, domain, validatorResult.Domain)
	})
}

func TestValidationDomainListMatchSetValidatorResultDomain(t *testing.T) {
	t.Run("validationDomainListMatch#setValidatorResultDomain", func(t *testing.T) {
		email, domain := pairRandomEmailDomain()
		validation := &validationDomainListMatch{result: createValidatorResult(email, createConfiguration())}
		validation.setValidatorResultDomain()

		assert.Equal(t, domain, validation.result.Domain)
	})
}

func TestValidationDomainListMatchIsWhitelistedDomain(t *testing.T) {
	email, domain := pairRandomEmailDomain()

	t.Run("when whitelisted domain", func(t *testing.T) {
		configuration, _ := NewConfiguration(ConfigurationAttr{VerifierEmail: randomEmail(), WhitelistedDomains: []string{domain}})
		validation := &validationDomainListMatch{result: createValidatorResult(email, configuration)}
		validation.setValidatorResultDomain()

		assert.True(t, validation.isWhitelistedDomain())
	})

	t.Run("when not whitelisted domain", func(t *testing.T) {
		validation := &validationDomainListMatch{result: createValidatorResult(email, createConfiguration())}
		validation.setValidatorResultDomain()

		assert.False(t, validation.isWhitelistedDomain())
	})
}

func TestValidationDomainListMatchIsWhitelistValidation(t *testing.T) {
	t.Run("when whitelist validation", func(t *testing.T) {
		configuration, _ := NewConfiguration(ConfigurationAttr{VerifierEmail: randomEmail(), WhitelistValidation: true})
		validation := &validationDomainListMatch{result: createValidatorResult(randomEmail(), configuration)}

		assert.True(t, validation.isWhitelistValidation())
	})

	t.Run("when not whitelist validation", func(t *testing.T) {
		validation := &validationDomainListMatch{result: createValidatorResult(randomEmail(), createConfiguration())}

		assert.False(t, validation.isWhitelistValidation())
	})
}

func TestValidationDomainListMatchIsBlacklistedDomain(t *testing.T) {
	email, domain := pairRandomEmailDomain()

	t.Run("when blacklisted domain", func(t *testing.T) {
		configuration, _ := NewConfiguration(ConfigurationAttr{VerifierEmail: randomEmail(), BlacklistedDomains: []string{domain}})
		validation := &validationDomainListMatch{result: createValidatorResult(email, configuration)}
		validation.setValidatorResultDomain()

		assert.True(t, validation.isBlacklistedDomain())
	})

	t.Run("when not blacklisted domain", func(t *testing.T) {
		validation := &validationDomainListMatch{result: createValidatorResult(email, createConfiguration())}
		validation.setValidatorResultDomain()

		assert.False(t, validation.isBlacklistedDomain())
	})
}
