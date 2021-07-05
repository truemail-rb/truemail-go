package truemail

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	// TODO: change to integration tests when validationMx.check() will be implemented
	for _, validValidationType := range []string{ValidationTypeRegex, ValidationTypeMx, ValidationTypeMxBlacklist, ValidationTypeSMTP} {
		t.Run(validValidationType+" valid validation type", func(t *testing.T) {
			email, configuration := randomEmail(), createConfiguration()
			validatorResult, err := Validate(email, configuration, validValidationType)

			assert.NoError(t, err)
			assert.Equal(t, email, validatorResult.Email)
			assert.Equal(t, configuration, validatorResult.Configuration)
			assert.Equal(t, validValidationType, validatorResult.ValidationType)
			assert.Equal(t, usedValidationsByType(validValidationType), validatorResult.usedValidations)
			assert.True(t, validatorResult.Success)
		})
	}

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
		configuration, _ := NewConfiguration(
			ConfigurationAttr{
				verifierEmail:       randomEmail(),
				whitelistedDomains:  []string{domain},
				whitelistValidation: true,
			},
		)
		validatorResult, _ := Validate(email, configuration)

		assert.True(t, validatorResult.Success)
		assert.True(t, validatorResult.isPassFromDomainListMatch)
		assert.Equal(t, usedValidationsByType(ValidationTypeDefault), validatorResult.usedValidations)
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
		email, ipAddress := randomEmail(), randomIpAddress()
		configuration, _ := NewConfiguration(
			ConfigurationAttr{
				verifierEmail:            randomEmail(),
				blacklistedMxIpAddresses: []string{ipAddress},
			},
		)
		validatorResult, _ := Validate(email, configuration, ValidationTypeMxBlacklist)

		// assert.False(t, validatorResult.Success) // TODO: update after validationMxBlacklist.check() implementation
		assert.Equal(t, usedValidationsByType(ValidationTypeMxBlacklist), validatorResult.usedValidations)
	})
}

func TestIsValid(t *testing.T) {
	t.Run("when succesful validation", func(t *testing.T) {
		assert.True(t, IsValid(randomEmail(), createConfiguration(), ValidationTypeRegex))
	})

	t.Run("when failure validation", func(t *testing.T) {
		assert.False(t, IsValid("invalid@email", createConfiguration(), ValidationTypeRegex))
	})

	t.Run("when invalid validation type", func(t *testing.T) {
		assert.False(t, IsValid(randomEmail(), createConfiguration(), "invalidValidationType"))
	})
}

func TestVariadicValidationType(t *testing.T) {
	t.Run("without validation type", func(t *testing.T) {
		result, err := variadicValidationType([]string{})

		assert.NoError(t, err)
		assert.Equal(t, ValidationTypeDefault, result)
	})

	t.Run("valid validation type", func(t *testing.T) {
		validationType := randomValidationType()
		result, err := variadicValidationType([]string{validationType})

		assert.NoError(t, err)
		assert.Equal(t, validationType, result)
	})

	t.Run("invalid validation type", func(t *testing.T) {
		invalidValidationType := "invalid type"
		result, err := variadicValidationType([]string{invalidValidationType})
		errorMessage := fmt.Sprintf("%s is invalid validation type, use one of these: [regex mx mx_blacklist smtp]", invalidValidationType)

		assert.EqualError(t, err, errorMessage)
		assert.Equal(t, invalidValidationType, result)
	})
}

func TestValidateValidationTypeContext(t *testing.T) {
	for _, validValidationType := range []string{ValidationTypeRegex, ValidationTypeMx, ValidationTypeSMTP} {
		t.Run("valid validation type", func(t *testing.T) {
			assert.NoError(t, validateValidationTypeContext(validValidationType))
		})
	}

	t.Run("invalid validation type", func(t *testing.T) {
		invalidType := "invalid type"
		errorMessage := fmt.Sprintf("%s is invalid validation type, use one of these: [regex mx mx_blacklist smtp]", invalidType)

		assert.EqualError(t, validateValidationTypeContext(invalidType), errorMessage)
	})
}

func TestNewValidator(t *testing.T) {
	t.Run("creates validator", func(t *testing.T) {
		email, validationType, configuration := randomEmail(), randomValidationType(), createConfiguration()
		validator := newValidator(email, validationType, configuration)
		validatorResult := validator.result

		assert.Equal(t, email, validatorResult.Email)
		assert.Equal(t, validationType, validatorResult.ValidationType)
		assert.Equal(t, configuration, validatorResult.Configuration)
		assert.False(t, validatorResult.isPassFromDomainListMatch)
		assert.Empty(t, validatorResult.usedValidations)
	})
}
