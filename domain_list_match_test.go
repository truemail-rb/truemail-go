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
				verifierEmail:      randomEmail(),
				whitelistedDomains: []string{domain},
			},
		)
		validatorResult := runDomainListMatchValidation(email, configuration)

		assert.True(t, validatorResult.Success)
		assert.False(t, validatorResult.isPassFromDomainListMatch)
		assert.Equal(t, DomainListMatchWhitelist, validatorResult.ValidationType)
	})

	t.Run("blacklist case, email is in blacklist", func(t *testing.T) {
		configuration, _ := NewConfiguration(
			ConfigurationAttr{
				verifierEmail:      randomEmail(),
				blacklistedDomains: []string{domain},
			},
		)
		validatorResult := runDomainListMatchValidation(email, configuration)

		assert.False(t, validatorResult.Success)
		assert.False(t, validatorResult.isPassFromDomainListMatch)
		assert.Equal(t, DomainListMatchBlacklist, validatorResult.ValidationType)
	})

	t.Run("whitelist/blackist case, email is not in both lists", func(t *testing.T) {
		configuration, _ := NewConfiguration(
			ConfigurationAttr{
				verifierEmail: randomEmail(),
			},
		)
		validationType := randomValidationType()
		validatorResult := runDomainListMatchValidation(email, configuration, validationType)

		assert.True(t, validatorResult.Success)
		assert.True(t, validatorResult.isPassFromDomainListMatch)
		assert.Equal(t, validationType, validatorResult.ValidationType)
	})

	t.Run("whitelist/blackist case, email is in both lists", func(t *testing.T) {
		configuration, _ := NewConfiguration(
			ConfigurationAttr{
				verifierEmail:      randomEmail(),
				whitelistedDomains: []string{domain},
				blacklistedDomains: []string{domain},
			},
		)
		validatorResult := runDomainListMatchValidation(email, configuration)

		assert.False(t, validatorResult.Success)
		assert.False(t, validatorResult.isPassFromDomainListMatch)
		assert.Equal(t, DomainListMatchBlacklist, validatorResult.ValidationType)
	})

	t.Run("whitelist validation case, email is in whitelist", func(t *testing.T) {
		configuration, _ := NewConfiguration(
			ConfigurationAttr{
				verifierEmail:       randomEmail(),
				whitelistValidation: true,
				whitelistedDomains:  []string{domain},
			},
		)
		validationType := randomValidationType()
		validatorResult := runDomainListMatchValidation(email, configuration, validationType)

		assert.True(t, validatorResult.Success)
		assert.True(t, validatorResult.isPassFromDomainListMatch)
		assert.Equal(t, validationType, validatorResult.ValidationType)
	})

	t.Run("whitelist validation case, email is not in whitelist/blacklist", func(t *testing.T) {
		configuration, _ := NewConfiguration(
			ConfigurationAttr{
				verifierEmail:       randomEmail(),
				whitelistValidation: true,
			},
		)
		validatorResult := runDomainListMatchValidation(email, configuration)

		assert.False(t, validatorResult.Success)
		assert.False(t, validatorResult.isPassFromDomainListMatch)
		assert.Equal(t, DomainListMatchBlacklist, validatorResult.ValidationType)
	})

	t.Run("whitelist validation case, email is in blacklist", func(t *testing.T) {
		configuration, _ := NewConfiguration(
			ConfigurationAttr{
				verifierEmail:       randomEmail(),
				whitelistValidation: true,
				blacklistedDomains:  []string{domain},
			},
		)
		validatorResult := runDomainListMatchValidation(email, configuration)

		assert.False(t, validatorResult.Success)
		assert.False(t, validatorResult.isPassFromDomainListMatch)
		assert.Equal(t, DomainListMatchBlacklist, validatorResult.ValidationType)
	})

	t.Run("whitelist validation case, email is in both lists", func(t *testing.T) {
		configuration, _ := NewConfiguration(
			ConfigurationAttr{
				verifierEmail:       randomEmail(),
				whitelistValidation: true,
				whitelistedDomains:  []string{domain},
				blacklistedDomains:  []string{domain},
			},
		)
		validatorResult := runDomainListMatchValidation(email, configuration)

		assert.False(t, validatorResult.Success)
		assert.False(t, validatorResult.isPassFromDomainListMatch)
		assert.Equal(t, DomainListMatchBlacklist, validatorResult.ValidationType)
	})
}

func TestValidationDomainListMatchIsWhitelistedDomain(t *testing.T) {
	email, domain := pairRandomEmailDomain()

	t.Run("when whitelisted domain", func(t *testing.T) {
		configuration, _ := NewConfiguration(ConfigurationAttr{verifierEmail: randomEmail(), whitelistedDomains: []string{domain}})
		validatorResult := createValidatorResult(email, configuration)

		assert.True(t, new(validationDomainListMatch).isWhitelistedDomain(validatorResult))
	})

	t.Run("when not whitelisted domain", func(t *testing.T) {
		validatorResult := createValidatorResult(email, createConfiguration())

		assert.False(t, new(validationDomainListMatch).isWhitelistedDomain(validatorResult))
	})
}

func TestValidationDomainListMatchIsWhitelistValidation(t *testing.T) {
	t.Run("when whitelist validation", func(t *testing.T) {
		configuration, _ := NewConfiguration(ConfigurationAttr{verifierEmail: randomEmail(), whitelistValidation: true})
		validatorResult := createValidatorResult(randomEmail(), configuration)

		assert.True(t, new(validationDomainListMatch).isWhitelistValidation(validatorResult))
	})

	t.Run("when not whitelist validation", func(t *testing.T) {
		configuration, _ := NewConfiguration(ConfigurationAttr{verifierEmail: randomEmail(), whitelistValidation: false})
		validatorResult := createValidatorResult(randomEmail(), configuration)

		assert.False(t, new(validationDomainListMatch).isWhitelistValidation(validatorResult))
	})
}

func TestValidationDomainListMatchIsBlacklistedDomain(t *testing.T) {
	email, domain := pairRandomEmailDomain()

	t.Run("when blacklisted domain", func(t *testing.T) {
		configuration, _ := NewConfiguration(ConfigurationAttr{verifierEmail: randomEmail(), blacklistedDomains: []string{domain}})
		validatorResult := createValidatorResult(email, configuration)

		assert.True(t, new(validationDomainListMatch).isBlacklistedDomain(validatorResult))
	})

	t.Run("when not blacklisted domain", func(t *testing.T) {
		validatorResult := createValidatorResult(email, createConfiguration())

		assert.False(t, new(validationDomainListMatch).isBlacklistedDomain(validatorResult))
	})
}
