package truemail

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	for _, validValidationType := range []string{"", ValidationTypeRegex, ValidationTypeMx, ValidationTypeSMTP} {
		t.Run(validValidationType+"valid validation type", func(t *testing.T) { // TODO: add stub for validator layers
			email, configuration := createRandomEmail(), createConfiguration()
			validationAttr := ValidationAttr{email: email, configuration: configuration, validationType: validValidationType}
			validatorResult, err := Validate(validationAttr)
			assert.NoError(t, err)
			assert.Equal(t, email, validatorResult.Email)
			assert.Equal(t, configuration, validatorResult.Configuration)
		})
	}

	t.Run("invalid validation type", func(t *testing.T) {
		invalidType := "invalid type"
		errorMessage := fmt.Sprintf("%s is invalid validation type, use one of these: [regex mx smtp]", invalidType)
		validationAttr := ValidationAttr{email: createRandomEmail(), configuration: createConfiguration(), validationType: invalidType}
		_, err := Validate(validationAttr)
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
		validatorResult, _ := Validate(ValidationAttr{email: email, configuration: configuration})
		assert.True(t, validatorResult.Success)
		assert.False(t, validatorResult.isPassFromDomainListMatch)
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
		validatorResult, _ := Validate(ValidationAttr{email: email, configuration: configuration})
		assert.True(t, validatorResult.Success)
		assert.True(t, validatorResult.isPassFromDomainListMatch)
	})

	t.Run("Whitelist/Blacklist validation fails", func(t *testing.T) {
		email, domain := createPairRandomEmailDomain()
		configuration, _ := NewConfiguration(
			ConfigurationAttr{
				verifierEmail:      createRandomEmail(),
				blacklistedDomains: []string{domain},
			},
		)
		validatorResult, _ := Validate(ValidationAttr{email: email, configuration: configuration})
		assert.False(t, validatorResult.Success)
		assert.False(t, validatorResult.isPassFromDomainListMatch)
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

func TestNewValidatorResult(t *testing.T) {
	t.Run("creates ValidatorResult", func(t *testing.T) {
		email, configuration := createRandomEmail(), createConfiguration()
		validatorResult := newValidatorResult(email, configuration, createRandomValidationType())
		assert.Equal(t, email, validatorResult.Email)
		assert.Equal(t, configuration, validatorResult.Configuration)
	})
}

func TestAddError(t *testing.T) {
	t.Run("addes error to ValidatorResult", func(t *testing.T) {
		key, value := "some_error_key", "some_error_value"
		validatorResult := addError(new(validatorResult), key, value)
		assert.Equal(t, value, validatorResult.Errors[key])
	})
}
