package truemail

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	// TODO: change to integration tests when .validateMx() will be implemented
	for _, validValidationType := range []string{ValidationTypeRegex, ValidationTypeMx, ValidationTypeSMTP} {
		t.Run(validValidationType+" valid validation type", func(t *testing.T) {
			email, configuration := createRandomEmail(), createConfiguration()
			validatorResult, err := Validate(email, configuration, validValidationType)
			assert.NoError(t, err)
			assert.Equal(t, email, validatorResult.Email)
			assert.Equal(t, configuration, validatorResult.Configuration)
			assert.Equal(t, validValidationType, validatorResult.ValidationType)
			assert.Equal(t, usedValidationsByType(validValidationType), validatorResult.validator.usedValidations)
			assert.True(t, validatorResult.Success)
		})
	}

	t.Run("invalid validation type", func(t *testing.T) {
		invalidValidationType := "invalid type"
		errorMessage := fmt.Sprintf("%s is invalid validation type, use one of these: [regex mx smtp]", invalidValidationType)
		_, err := Validate(createRandomEmail(), createConfiguration(), invalidValidationType)
		assert.EqualError(t, err, errorMessage)
	})

	t.Run("Whitelist/Blacklist validation successful", func(t *testing.T) {
		email, domain := createPairRandomEmailDomain()
		configuration, _ := NewConfiguration(
			ConfigurationAttr{
				verifierEmail:      createRandomEmail(),
				whitelistedDomains: []string{domain},
			},
		)
		validatorResult, _ := Validate(email, configuration)
		assert.True(t, validatorResult.Success)
		assert.False(t, validatorResult.validator.isPassFromDomainListMatch)
		assert.Empty(t, validatorResult.validator.usedValidations)
	})

	t.Run("Whitelist/Blacklist validation passes to next validation level", func(t *testing.T) {
		email, domain := createPairRandomEmailDomain()
		configuration, _ := NewConfiguration(
			ConfigurationAttr{
				verifierEmail:       createRandomEmail(),
				whitelistedDomains:  []string{domain},
				whitelistValidation: true,
			},
		)
		validatorResult, _ := Validate(email, configuration)
		assert.True(t, validatorResult.Success)
		assert.True(t, validatorResult.validator.isPassFromDomainListMatch)
		assert.Equal(t, usedValidationsByType(ValidationTypeDefault), validatorResult.validator.usedValidations)
	})

	t.Run("Whitelist/Blacklist validation fails", func(t *testing.T) {
		email, domain := createPairRandomEmailDomain()
		configuration, _ := NewConfiguration(
			ConfigurationAttr{
				verifierEmail:      createRandomEmail(),
				blacklistedDomains: []string{domain},
			},
		)
		validatorResult, _ := Validate(email, configuration)
		assert.False(t, validatorResult.Success)
		assert.False(t, validatorResult.validator.isPassFromDomainListMatch)
		assert.Empty(t, validatorResult.validator.usedValidations)
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
		errorMessage := fmt.Sprintf("%s is invalid validation type, use one of these: [regex mx smtp]", invalidType)
		assert.EqualError(t, validateValidationTypeContext(invalidType), errorMessage)
	})
}

func TestNewValidator(t *testing.T) {
	t.Run("creates validator", func(t *testing.T) {
		email, validationType, configuration := createRandomEmail(), createRandomValidationType(), createConfiguration()
		validator := newValidator(email, validationType, configuration)
		validatorResult := validator.result
		assert.Equal(t, email, validatorResult.Email)
		assert.Equal(t, validationType, validatorResult.ValidationType)
		assert.Equal(t, configuration, validatorResult.Configuration)
		assert.Equal(t, validator, validatorResult.validator)
	})
}

func TestAddError(t *testing.T) {
	t.Run("addes error to ValidatorResult", func(t *testing.T) {
		key, value := "some_error_key", "some_error_value"
		validatorResult := addError(new(validatorResult), key, value)
		assert.Equal(t, value, validatorResult.Errors[key])
	})
}
