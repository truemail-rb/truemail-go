package truemail

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateDomainListMatch(t *testing.T) {
	email, domain := createPairRandomEmailDomain()

	t.Run("whitelist case, email is in whitelist", func(t *testing.T) {
		configuration, _ := NewConfiguration(
			ConfigurationAttr{
				verifierEmail:      createRandomEmail(),
				whitelistedDomains: []string{domain},
			},
		)
		validatorResult := validateDomainListMatch(&validatorResult{Email: email, Configuration: configuration})
		assert.True(t, validatorResult.Success)
		assert.False(t, validatorResult.isPassFromDomainListMatch)
		assert.Equal(t, DomainListMatchWhitelist, validatorResult.ValidationType)
	})

	t.Run("blacklist case, email is in blacklist", func(t *testing.T) {
		configuration, _ := NewConfiguration(
			ConfigurationAttr{
				verifierEmail:      createRandomEmail(),
				blacklistedDomains: []string{domain},
			},
		)
		validatorResult := validateDomainListMatch(&validatorResult{Email: email, Configuration: configuration})
		assert.False(t, validatorResult.Success)
		assert.False(t, validatorResult.isPassFromDomainListMatch)
		assert.Equal(t, DomainListMatchBlacklist, validatorResult.ValidationType)
	})

	t.Run("whitelist/blackist case, email is not in both lists", func(t *testing.T) {
		configuration, _ := NewConfiguration(
			ConfigurationAttr{
				verifierEmail: createRandomEmail(),
			},
		)
		validationType := createRandomValidationType()
		validatorResult := validateDomainListMatch(
			&validatorResult{
				Email:          email,
				Configuration:  configuration,
				ValidationType: validationType,
			},
		)
		assert.True(t, validatorResult.Success)
		assert.True(t, validatorResult.isPassFromDomainListMatch)
		assert.Equal(t, validationType, validatorResult.ValidationType)
	})

	t.Run("whitelist/blackist case, email is in both lists", func(t *testing.T) {
		configuration, _ := NewConfiguration(
			ConfigurationAttr{
				verifierEmail:      createRandomEmail(),
				whitelistedDomains: []string{domain},
				blacklistedDomains: []string{domain},
			},
		)
		validatorResult := validateDomainListMatch(&validatorResult{Email: email, Configuration: configuration})
		assert.False(t, validatorResult.Success)
		assert.False(t, validatorResult.isPassFromDomainListMatch)
		assert.Equal(t, DomainListMatchBlacklist, validatorResult.ValidationType)
	})

	t.Run("whitelist validation case, email is in whitelist", func(t *testing.T) {
		configuration, _ := NewConfiguration(
			ConfigurationAttr{
				verifierEmail:       createRandomEmail(),
				whitelistValidation: true,
				whitelistedDomains:  []string{domain},
			},
		)
		validationType := createRandomValidationType()
		validatorResult := validateDomainListMatch(
			&validatorResult{
				Email:          email,
				Configuration:  configuration,
				ValidationType: validationType,
			},
		)
		assert.True(t, validatorResult.Success)
		assert.True(t, validatorResult.isPassFromDomainListMatch)
		assert.Equal(t, validationType, validatorResult.ValidationType)
	})

	t.Run("whitelist validation case, email is not in whitelist/blacklist", func(t *testing.T) {
		configuration, _ := NewConfiguration(
			ConfigurationAttr{
				verifierEmail:       createRandomEmail(),
				whitelistValidation: true,
			},
		)
		validatorResult := validateDomainListMatch(&validatorResult{Email: email, Configuration: configuration})
		assert.False(t, validatorResult.Success)
		assert.False(t, validatorResult.isPassFromDomainListMatch)
		assert.Equal(t, DomainListMatchBlacklist, validatorResult.ValidationType)
	})

	t.Run("whitelist validation case, email is in blacklist", func(t *testing.T) {
		configuration, _ := NewConfiguration(
			ConfigurationAttr{
				verifierEmail:       createRandomEmail(),
				whitelistValidation: true,
				blacklistedDomains:  []string{domain},
			},
		)
		validatorResult := validateDomainListMatch(&validatorResult{Email: email, Configuration: configuration})
		assert.False(t, validatorResult.Success)
		assert.False(t, validatorResult.isPassFromDomainListMatch)
		assert.Equal(t, DomainListMatchBlacklist, validatorResult.ValidationType)
	})

	t.Run("whitelist validation case, email is in both lists", func(t *testing.T) {
		configuration, _ := NewConfiguration(
			ConfigurationAttr{
				verifierEmail:       createRandomEmail(),
				whitelistValidation: true,
				whitelistedDomains:  []string{domain},
				blacklistedDomains:  []string{domain},
			},
		)
		validatorResult := validateDomainListMatch(&validatorResult{Email: email, Configuration: configuration})
		assert.False(t, validatorResult.Success)
		assert.False(t, validatorResult.isPassFromDomainListMatch)
		assert.Equal(t, DomainListMatchBlacklist, validatorResult.ValidationType)
	})
}

func TestEmailDomain(t *testing.T) {
	t.Run("extracts domain name from email address", func(t *testing.T) {
		email, domain := createPairRandomEmailDomain()
		assert.Equal(t, domain, emailDomain(email))
	})
}

func TestIsWhitelistedDomain(t *testing.T) {
	email, domain := createPairRandomEmailDomain()

	t.Run("when whitelisted domain", func(t *testing.T) {
		configuration, _ := NewConfiguration(ConfigurationAttr{verifierEmail: createRandomEmail(), whitelistedDomains: []string{domain}})
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
		configuration, _ := NewConfiguration(ConfigurationAttr{verifierEmail: createRandomEmail(), whitelistValidation: true})
		validatorResult := createValidatorResult(createRandomEmail(), configuration)
		assert.True(t, isWhitelistValidation(validatorResult))
	})

	t.Run("when not whitelist validation", func(t *testing.T) {
		configuration, _ := NewConfiguration(ConfigurationAttr{verifierEmail: createRandomEmail(), whitelistValidation: false})
		validatorResult := createValidatorResult(createRandomEmail(), configuration)
		assert.False(t, isWhitelistValidation(validatorResult))
	})
}

func TestIsBlacklistedDomain(t *testing.T) {
	email, domain := createPairRandomEmailDomain()

	t.Run("when blacklisted domain", func(t *testing.T) {
		configuration, _ := NewConfiguration(ConfigurationAttr{verifierEmail: createRandomEmail(), blacklistedDomains: []string{domain}})
		validatorResult := createValidatorResult(email, configuration)
		assert.True(t, isBlacklistedDomain(validatorResult))
	})

	t.Run("when not blacklisted domain", func(t *testing.T) {
		validatorResult := createValidatorResult(email, createConfiguration())
		assert.False(t, isBlacklistedDomain(validatorResult))
	})
}
