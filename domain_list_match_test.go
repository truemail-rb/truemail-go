package truemail

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateDomainListMatch(t *testing.T) {
	email, domain := pairRandomEmailDomain()

	t.Run("whitelist case, email is in whitelist", func(t *testing.T) {
		configuration, _ := NewConfiguration(
			ConfigurationAttr{
				verifierEmail:      randomEmail(),
				whitelistedDomains: []string{domain},
			},
		)
		validatorResult := validateDomainListMatch(createValidatorResult(email, configuration))
		assert.True(t, validatorResult.Success)
		assert.False(t, validatorResult.validator.isPassFromDomainListMatch)
		assert.Equal(t, DomainListMatchWhitelist, validatorResult.ValidationType)
	})

	t.Run("blacklist case, email is in blacklist", func(t *testing.T) {
		configuration, _ := NewConfiguration(
			ConfigurationAttr{
				verifierEmail:      randomEmail(),
				blacklistedDomains: []string{domain},
			},
		)
		validatorResult := validateDomainListMatch(createValidatorResult(email, configuration))
		assert.False(t, validatorResult.Success)
		assert.False(t, validatorResult.validator.isPassFromDomainListMatch)
		assert.Equal(t, DomainListMatchBlacklist, validatorResult.ValidationType)
	})

	t.Run("whitelist/blackist case, email is not in both lists", func(t *testing.T) {
		configuration, _ := NewConfiguration(
			ConfigurationAttr{
				verifierEmail: randomEmail(),
			},
		)
		validationType := createRandomValidationType()
		validatorResult := validateDomainListMatch(createValidatorResult(email, configuration, validationType))

		assert.True(t, validatorResult.Success)
		assert.True(t, validatorResult.validator.isPassFromDomainListMatch)
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
		validatorResult := validateDomainListMatch(createValidatorResult(email, configuration))
		assert.False(t, validatorResult.Success)
		assert.False(t, validatorResult.validator.isPassFromDomainListMatch)
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
		validationType := createRandomValidationType()
		validatorResult := validateDomainListMatch(createValidatorResult(email, configuration, validationType))
		assert.True(t, validatorResult.Success)
		assert.True(t, validatorResult.validator.isPassFromDomainListMatch)
		assert.Equal(t, validationType, validatorResult.ValidationType)
	})

	t.Run("whitelist validation case, email is not in whitelist/blacklist", func(t *testing.T) {
		configuration, _ := NewConfiguration(
			ConfigurationAttr{
				verifierEmail:       randomEmail(),
				whitelistValidation: true,
			},
		)
		validatorResult := validateDomainListMatch(createValidatorResult(email, configuration))
		assert.False(t, validatorResult.Success)
		assert.False(t, validatorResult.validator.isPassFromDomainListMatch)
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
		validatorResult := validateDomainListMatch(createValidatorResult(email, configuration))
		assert.False(t, validatorResult.Success)
		assert.False(t, validatorResult.validator.isPassFromDomainListMatch)
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
		validatorResult := validateDomainListMatch(createValidatorResult(email, configuration))
		assert.False(t, validatorResult.Success)
		assert.False(t, validatorResult.validator.isPassFromDomainListMatch)
		assert.Equal(t, DomainListMatchBlacklist, validatorResult.ValidationType)
	})
}

func TestEmailDomain(t *testing.T) {
	t.Run("extracts domain name from email address", func(t *testing.T) {
		email, domain := pairRandomEmailDomain()
		assert.Equal(t, domain, emailDomain(email))
	})
}

func TestIsWhitelistedDomain(t *testing.T) {
	email, domain := pairRandomEmailDomain()

	t.Run("when whitelisted domain", func(t *testing.T) {
		configuration, _ := NewConfiguration(ConfigurationAttr{verifierEmail: randomEmail(), whitelistedDomains: []string{domain}})
		validatorResult := createValidatorResult(email, configuration)
		assert.True(t, isWhitelistedDomain(validatorResult))
	})

	t.Run("when not whitelisted domain", func(t *testing.T) {
		validatorResult := createValidatorResult(email, createConfiguration())
		assert.False(t, isWhitelistedDomain(validatorResult))
	})
}

func TestIsWhitelistValidation(t *testing.T) {
	t.Run("when whitelist validation", func(t *testing.T) {
		configuration, _ := NewConfiguration(ConfigurationAttr{verifierEmail: randomEmail(), whitelistValidation: true})
		validatorResult := createValidatorResult(randomEmail(), configuration)
		assert.True(t, isWhitelistValidation(validatorResult))
	})

	t.Run("when not whitelist validation", func(t *testing.T) {
		configuration, _ := NewConfiguration(ConfigurationAttr{verifierEmail: randomEmail(), whitelistValidation: false})
		validatorResult := createValidatorResult(randomEmail(), configuration)
		assert.False(t, isWhitelistValidation(validatorResult))
	})
}

func TestIsBlacklistedDomain(t *testing.T) {
	email, domain := pairRandomEmailDomain()

	t.Run("when blacklisted domain", func(t *testing.T) {
		configuration, _ := NewConfiguration(ConfigurationAttr{verifierEmail: randomEmail(), blacklistedDomains: []string{domain}})
		validatorResult := createValidatorResult(email, configuration)
		assert.True(t, isBlacklistedDomain(validatorResult))
	})

	t.Run("when not blacklisted domain", func(t *testing.T) {
		validatorResult := createValidatorResult(email, createConfiguration())
		assert.False(t, isBlacklistedDomain(validatorResult))
	})
}
