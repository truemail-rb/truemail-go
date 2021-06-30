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
			email, configuration := randomEmail(), createConfiguration()
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
		assert.False(t, validatorResult.validator.isPassFromDomainListMatch)
		assert.Empty(t, validatorResult.validator.usedValidations)
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
		assert.True(t, validatorResult.validator.isPassFromDomainListMatch)
		assert.Equal(t, usedValidationsByType(ValidationTypeDefault), validatorResult.validator.usedValidations)
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
		assert.False(t, validatorResult.validator.isPassFromDomainListMatch)
		assert.Empty(t, validatorResult.validator.usedValidations)
	})

	// t.Run("Mx blacklist validation fails", func(t *testing.T) { // TODO: update after layer implementation
	// 	email, ipAddress := randomEmail(), randomIpAddress()
	// 	configuration, _ := NewConfiguration(
	// 		ConfigurationAttr{
	// 			verifierEmail:            randomEmail(),
	// 			blacklistedMxIpAddresses: []string{ipAddress},
	// 		},
	// 	)
	// 	validatorResult, _ := Validate(email, configuration)
	//
	// 	assert.False(t, validatorResult.Success)
	// 	assert.Equal(t, usedValidationsByType(ValidationTypeMx), validatorResult.validator.usedValidations)
	// })
}

func TestValidatorValidateDomainListMatch(t *testing.T) {
	t.Run("#validateDomainListMatch", func(t *testing.T) {
		validator := createValidatorPassedFromDomainListMatch(randomEmail(), createConfiguration(), ValidationTypeRegex)
		validation, result := new(validationMock), validator.result
		validator.validate = validation

		validation.On("domainListMatch", result).Return(result)
		validator.validateDomainListMatch()
		validation.AssertExpectations(t)
	})
}

func TestValidatorValidateRegex(t *testing.T) {
	t.Run("#validateRegex", func(t *testing.T) {
		validator := createValidatorPassedFromDomainListMatch(randomEmail(), createConfiguration(), ValidationTypeRegex)
		validation, result := new(validationMock), validator.result
		validator.validate = validation

		validation.On(ValidationTypeRegex, result).Return(result)
		validator.validateRegex()
		validation.AssertExpectations(t)
	})
}
func TestValidatorValidateMx(t *testing.T) {
	t.Run("#validateMx", func(t *testing.T) {
		validator := createValidatorPassedFromDomainListMatch(randomEmail(), createConfiguration(), ValidationTypeRegex)
		validation, result := new(validationMock), validator.result
		validator.validate = validation

		validation.On(ValidationTypeMx, result).Return(result)
		validator.validateMx()
		validation.AssertExpectations(t)
	})
}

func TestValidatorValidateSMTP(t *testing.T) {
	t.Run("#validateSMTP", func(t *testing.T) {
		validator := createValidatorPassedFromDomainListMatch(randomEmail(), createConfiguration(), ValidationTypeRegex)
		validation, result := new(validationMock), validator.result
		validator.validate = validation

		validation.On(ValidationTypeSMTP, result).Return(result)
		validator.validateSMTP()
		validation.AssertExpectations(t)
	})
}

func TestValidatorRun(t *testing.T) {
	t.Run("when domainListMatch not passes", func(t *testing.T) {
		validator := createValidator(randomEmail(), createConfiguration(), ValidationTypeRegex)
		validation, result := new(validationMock), validator.result
		validator.validate = validation

		validation.On("domainListMatch", result).Return(result)
		assert.Equal(t, result, validator.run())
		validation.AssertExpectations(t)
	})

	t.Run("when domainListMatch passes, regex", func(t *testing.T) {
		validator := createValidatorPassedFromDomainListMatch(randomEmail(), createConfiguration(), ValidationTypeRegex)
		validation, result := new(validationMock), validator.result
		validator.validate = validation

		validation.On("domainListMatch", result).Return(result)
		validation.On(ValidationTypeRegex, result).Return(result)
		assert.Equal(t, result, validator.run())
		validation.AssertExpectations(t)
	})

	t.Run("when domainListMatch passes, mx", func(t *testing.T) {
		validator := createValidatorPassedFromDomainListMatch(randomEmail(), createConfiguration(), ValidationTypeMx)
		validation, result := new(validationMock), validator.result
		validator.validate = validation

		validation.On("domainListMatch", result).Return(result)
		validation.On(ValidationTypeMx, result).Return(result)
		assert.Equal(t, result, validator.run())
		validation.AssertExpectations(t)
	})

	t.Run("when domainListMatch passes, smtp", func(t *testing.T) {
		validator := createValidatorPassedFromDomainListMatch(randomEmail(), createConfiguration(), ValidationTypeSMTP)
		validation, result := new(validationMock), validator.result
		validator.validate = validation

		validation.On("domainListMatch", result).Return(result)
		validation.On(ValidationTypeSMTP, result).Return(result)
		assert.Equal(t, result, validator.run())
		validation.AssertExpectations(t)
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
		errorMessage := fmt.Sprintf("%s is invalid validation type, use one of these: [regex mx smtp]", invalidValidationType)

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
		errorMessage := fmt.Sprintf("%s is invalid validation type, use one of these: [regex mx smtp]", invalidType)

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
