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
	// 	assert.Equal(t, usedValidationsByType(ValidationTypeMx), validatorResult.usedValidations)
	// })
}

func TestValidatorValidateDomainListMatch(t *testing.T) {
	t.Run("#validateDomainListMatch", func(t *testing.T) {
		validator := createValidator(randomEmail(), createConfiguration())
		validationDomainListMatch, result := new(validationDomainListMatchMock), validator.result
		validator.domainListMatch = validationDomainListMatch

		validationDomainListMatch.On("check", result).Return(result)
		validator.validateDomainListMatch()
		validationDomainListMatch.AssertExpectations(t)
	})
}

func TestValidatorValidateRegex(t *testing.T) {
	t.Run("#validateRegex", func(t *testing.T) {
		validator := createValidator(randomEmail(), createConfiguration())
		validationRegex, result := new(validationRegexMock), validator.result
		validator.regex = validationRegex

		validationRegex.On("check", result).Return(result)
		validator.validateRegex()
		validationRegex.AssertExpectations(t)
		assert.Equal(t, usedValidationsByType(ValidationTypeRegex), validator.result.usedValidations)

	})
}
func TestValidatorValidateMx(t *testing.T) {
	t.Run("when all layers passed", func(t *testing.T) {
		validator := createValidator(randomEmail(), createConfiguration())
		validationRegex, validationMx, result := new(validationRegexMock), new(validationMxMock), validator.result
		validator.regex, validator.mx = validationRegex, validationMx
		result.Success = true

		validationRegex.On("check", result).Return(result)
		validationMx.On("check", result).Return(result)
		validator.validateMx()
		validationRegex.AssertExpectations(t)
		validationMx.AssertExpectations(t)
		assert.Equal(t, usedValidationsByType(ValidationTypeMx), validator.result.usedValidations)
	})

	t.Run("when regex layer fails", func(t *testing.T) {
		validator := createValidator(randomEmail(), createConfiguration())
		validationRegex, validationMx, result := new(validationRegexMock), new(validationMxMock), validator.result
		validator.regex, validator.mx = validationRegex, validationMx

		validationRegex.On("check", result).Return(result)
		validator.validateMx()
		validationRegex.AssertExpectations(t)
		validationMx.AssertNotCalled(t, "check", result)
		assert.Equal(t, usedValidationsByType(ValidationTypeRegex), validator.result.usedValidations)
	})
}

func TestValidatorValidateMxBlacklist(t *testing.T) {
	t.Run("when all layers passed", func(t *testing.T) {
		validator := createValidator(randomEmail(), createConfiguration())
		validationRegex, validationMx, validationMxBlacklist, result := new(validationRegexMock), new(validationMxMock), new(validationMxBlacklistMock), validator.result
		validator.regex, validator.mx, validator.mxBlacklist = validationRegex, validationMx, validationMxBlacklist
		result.Success = true

		validationRegex.On("check", result).Return(result)
		validationMx.On("check", result).Return(result)
		validationMxBlacklist.On("check", result).Return(result)
		validator.validateMxBlacklist()
		validationRegex.AssertExpectations(t)
		validationMx.AssertExpectations(t)
		validationMxBlacklist.AssertExpectations(t)
		assert.Equal(t, usedValidationsByType(ValidationTypeMxBlacklist), validator.result.usedValidations)
	})

	t.Run("when regex layer fails", func(t *testing.T) {
		validator := createValidator(randomEmail(), createConfiguration())
		validationRegex, validationMx, validationMxBlacklist, result := new(validationRegexMock), new(validationMxMock), new(validationMxBlacklistMock), validator.result
		validator.regex, validator.mx, validator.mxBlacklist = validationRegex, validationMx, validationMxBlacklist

		validationRegex.On("check", result).Return(result)
		validator.validateMxBlacklist()
		validationRegex.AssertExpectations(t)
		validationMx.AssertNotCalled(t, "check", result)
		validationMxBlacklist.AssertNotCalled(t, "check", result)
		assert.Equal(t, usedValidationsByType(ValidationTypeRegex), validator.result.usedValidations)
	})

	t.Run("when mx layer fails", func(t *testing.T) {
		validator := createValidator(randomEmail(), createConfiguration())
		validationRegex, validationMx, validationMxBlacklist, result := new(validationRegexMock), new(validationMxMock), new(validationMxBlacklistMock), validator.result
		validator.regex, validator.mx, validator.mxBlacklist = validationRegex, validationMx, validationMxBlacklist
		result.Success = true
		failedResult := failedValidatorResult()

		validationRegex.On("check", result).Return(result)
		validationMx.On("check", result).Return(failedResult)
		validator.validateMxBlacklist()
		validationRegex.AssertExpectations(t)
		validationMx.AssertExpectations(t)
		validationMxBlacklist.AssertNotCalled(t, "check", failedResult)
		assert.Equal(t, usedValidationsByType(ValidationTypeMx), validator.result.usedValidations)
	})
}

func TestValidatorValidateSMTP(t *testing.T) {
	t.Run("when all layers passed", func(t *testing.T) {
		validator := createValidator(randomEmail(), createConfiguration())
		validationRegex, validationMx, validationMxBlacklist, validationSmtp := new(validationRegexMock), new(validationMxMock), new(validationMxBlacklistMock), new(validationSmtpMock)
		validator.regex, validator.mx, validator.mxBlacklist, validator.smtp = validationRegex, validationMx, validationMxBlacklist, validationSmtp
		result := validator.result
		result.Success = true

		validationRegex.On("check", result).Return(result)
		validationMx.On("check", result).Return(result)
		validationMxBlacklist.On("check", result).Return(result)
		validationSmtp.On("check", result).Return(result)
		validator.validateSMTP()
		validationRegex.AssertExpectations(t)
		validationMx.AssertExpectations(t)
		validationMxBlacklist.AssertExpectations(t)
		validationSmtp.AssertExpectations(t)
		assert.Equal(t, usedValidationsByType(ValidationTypeSMTP), validator.result.usedValidations)
	})

	t.Run("when regex layer fails", func(t *testing.T) {
		validator := createValidator(randomEmail(), createConfiguration())
		validationRegex, validationMx, validationMxBlacklist, validationSmtp := new(validationRegexMock), new(validationMxMock), new(validationMxBlacklistMock), new(validationSmtpMock)
		validator.regex, validator.mx, validator.mxBlacklist, validator.smtp = validationRegex, validationMx, validationMxBlacklist, validationSmtp
		result := validator.result
		result.Success = true
		failedResult := failedValidatorResult()

		validationRegex.On("check", result).Return(failedResult)
		validator.validateSMTP()
		validationRegex.AssertExpectations(t)
		validationMx.AssertNotCalled(t, "check", failedResult)
		validationMxBlacklist.AssertNotCalled(t, "check", failedResult)
		validationSmtp.AssertNotCalled(t, "check", failedResult)
		assert.Equal(t, usedValidationsByType(ValidationTypeRegex), validator.result.usedValidations)
	})

	t.Run("when mx layer fails", func(t *testing.T) {
		validator := createValidator(randomEmail(), createConfiguration())
		validationRegex, validationMx, validationMxBlacklist, validationSmtp := new(validationRegexMock), new(validationMxMock), new(validationMxBlacklistMock), new(validationSmtpMock)
		validator.regex, validator.mx, validator.mxBlacklist, validator.smtp = validationRegex, validationMx, validationMxBlacklist, validationSmtp
		result := validator.result
		result.Success = true
		failedResult := failedValidatorResult()

		validationRegex.On("check", result).Return(result)
		validationMx.On("check", result).Return(failedResult)
		validator.validateSMTP()
		validationRegex.AssertExpectations(t)
		validationMx.AssertExpectations(t)
		validationMxBlacklist.AssertNotCalled(t, "check", failedResult)
		validationSmtp.AssertNotCalled(t, "check", failedResult)
		assert.Equal(t, usedValidationsByType(ValidationTypeMx), validator.result.usedValidations)
	})

	t.Run("when mx blacklist layer fails", func(t *testing.T) {
		validator := createValidator(randomEmail(), createConfiguration())
		validationRegex, validationMx, validationMxBlacklist, validationSmtp := new(validationRegexMock), new(validationMxMock), new(validationMxBlacklistMock), new(validationSmtpMock)
		validator.regex, validator.mx, validator.mxBlacklist, validator.smtp = validationRegex, validationMx, validationMxBlacklist, validationSmtp
		result := validator.result
		result.Success = true
		failedResult := failedValidatorResult()

		validationRegex.On("check", result).Return(result)
		validationMx.On("check", result).Return(result)
		validationMxBlacklist.On("check", result).Return(failedResult)
		validator.validateSMTP()
		validationRegex.AssertExpectations(t)
		validationMx.AssertExpectations(t)
		validationMxBlacklist.AssertExpectations(t)
		validationSmtp.AssertNotCalled(t, "check", failedResult)
		assert.Equal(t, usedValidationsByType(ValidationTypeMxBlacklist), validator.result.usedValidations)
	})

	t.Run("when smtp layer fails", func(t *testing.T) {
		validator := createValidator(randomEmail(), createConfiguration())
		validationRegex, validationMx, validationMxBlacklist, validationSmtp := new(validationRegexMock), new(validationMxMock), new(validationMxBlacklistMock), new(validationSmtpMock)
		validator.regex, validator.mx, validator.mxBlacklist, validator.smtp = validationRegex, validationMx, validationMxBlacklist, validationSmtp
		result := validator.result
		result.Success = true
		failedResult := failedValidatorResult()

		validationRegex.On("check", result).Return(result)
		validationMx.On("check", result).Return(result)
		validationMxBlacklist.On("check", result).Return(result)
		validationSmtp.On("check", result).Return(failedResult)
		validator.validateSMTP()
		validationRegex.AssertExpectations(t)
		validationMx.AssertExpectations(t)
		validationMxBlacklist.AssertExpectations(t)
		validationSmtp.AssertExpectations(t)
		assert.Equal(t, usedValidationsByType(ValidationTypeSMTP), validator.result.usedValidations)
	})
}

func TestValidatorRun(t *testing.T) {
	t.Run("domainListMatch fails", func(t *testing.T) {
		validator := createValidator(randomEmail(), createConfiguration(), ValidationTypeRegex)
		validationDomainListMatch, result := new(validationDomainListMatchMock), validator.result
		validationRegex, validationMx, validationMxBlacklist, validationSmtp := new(validationRegexMock), new(validationMxMock), new(validationMxBlacklistMock), new(validationSmtpMock)
		validator.domainListMatch, validator.regex, validator.mx, validator.mxBlacklist, validator.smtp = validationDomainListMatch, validationRegex, validationMx, validationMxBlacklist, validationSmtp

		validationDomainListMatch.On("check", result).Return(result)
		assert.Equal(t, result, validator.run())
		validationDomainListMatch.AssertExpectations(t)
		validationRegex.AssertNotCalled(t, "check", result)
		validationMx.AssertNotCalled(t, "check", result)
		validationMxBlacklist.AssertNotCalled(t, "check", result)
		validationSmtp.AssertNotCalled(t, "check", result)
	})

	t.Run("regex validation: domainListMatch succeed", func(t *testing.T) {
		validator := createValidator(randomEmail(), createConfiguration(), ValidationTypeRegex)
		validationDomainListMatch, result := new(validationDomainListMatchMock), validator.result
		validationRegex, validationMx, validationMxBlacklist, validationSmtp := new(validationRegexMock), new(validationMxMock), new(validationMxBlacklistMock), new(validationSmtpMock)
		validator.domainListMatch, validator.regex, validator.mx, validator.mxBlacklist, validator.smtp = validationDomainListMatch, validationRegex, validationMx, validationMxBlacklist, validationSmtp
		doPassedFromDomainListMatch(result)

		validationDomainListMatch.On("check", result).Return(result)
		validationRegex.On("check", result).Return(result)
		assert.Equal(t, result, validator.run())
		validationDomainListMatch.AssertExpectations(t)
		validationRegex.AssertExpectations(t)
		validationMx.AssertNotCalled(t, "check", result)
		validationMxBlacklist.AssertNotCalled(t, "check", result)
		validationSmtp.AssertNotCalled(t, "check", result)
	})

	t.Run("mx validation: domainListMatch succeed, regex fails", func(t *testing.T) {
		validator := createValidator(randomEmail(), createConfiguration(), ValidationTypeMx)
		validationDomainListMatch, result := new(validationDomainListMatchMock), validator.result
		validationRegex, validationMx, validationMxBlacklist, validationSmtp := new(validationRegexMock), new(validationMxMock), new(validationMxBlacklistMock), new(validationSmtpMock)
		validator.domainListMatch, validator.regex, validator.mx, validator.mxBlacklist, validator.smtp = validationDomainListMatch, validationRegex, validationMx, validationMxBlacklist, validationSmtp
		doPassedFromDomainListMatch(result)
		failedResult := failedValidatorResult()

		validationDomainListMatch.On("check", result).Return(result)
		validationRegex.On("check", result).Return(failedResult)
		validator.run()
		validationDomainListMatch.AssertExpectations(t)
		validationRegex.AssertExpectations(t)
		validationMx.AssertNotCalled(t, "check", failedResult)
		validationMxBlacklist.AssertNotCalled(t, "check", failedResult)
		validationSmtp.AssertNotCalled(t, "check", failedResult)
	})

	t.Run("mx validation: domainListMatch, regex succeed", func(t *testing.T) {
		validator := createValidator(randomEmail(), createConfiguration(), ValidationTypeMx)
		validationDomainListMatch, result := new(validationDomainListMatchMock), validator.result
		validationRegex, validationMx, validationMxBlacklist, validationSmtp := new(validationRegexMock), new(validationMxMock), new(validationMxBlacklistMock), new(validationSmtpMock)
		validator.domainListMatch, validator.regex, validator.mx, validator.mxBlacklist, validator.smtp = validationDomainListMatch, validationRegex, validationMx, validationMxBlacklist, validationSmtp
		doPassedFromDomainListMatch(result)

		validationDomainListMatch.On("check", result).Return(result)
		validationRegex.On("check", result).Return(result)
		validationMx.On("check", result).Return(result)
		assert.Equal(t, result, validator.run())
		validationDomainListMatch.AssertExpectations(t)
		validationRegex.AssertExpectations(t)
		validationMxBlacklist.AssertNotCalled(t, "check", result)
		validationSmtp.AssertNotCalled(t, "check", result)
	})

	t.Run("mx blacklist validation: domainListMatch succeed, regex fails", func(t *testing.T) {
		validator := createValidator(randomEmail(), createConfiguration(), ValidationTypeMxBlacklist)
		validationDomainListMatch, result := new(validationDomainListMatchMock), validator.result
		validationRegex, validationMx, validationMxBlacklist, validationSmtp := new(validationRegexMock), new(validationMxMock), new(validationMxBlacklistMock), new(validationSmtpMock)
		validator.domainListMatch, validator.regex, validator.mx, validator.mxBlacklist, validator.smtp = validationDomainListMatch, validationRegex, validationMx, validationMxBlacklist, validationSmtp
		doPassedFromDomainListMatch(result)
		failedResult := failedValidatorResult()

		validationDomainListMatch.On("check", result).Return(result)
		validationRegex.On("check", result).Return(failedResult)
		validator.run()
		validationDomainListMatch.AssertExpectations(t)
		validationRegex.AssertExpectations(t)
		validationMx.AssertNotCalled(t, "check", failedResult)
		validationMxBlacklist.AssertNotCalled(t, "check", failedResult)
		validationSmtp.AssertNotCalled(t, "check", failedResult)
	})

	t.Run("mx blacklist validation: domainListMatch, regex succeed, mx fails", func(t *testing.T) {
		validator := createValidator(randomEmail(), createConfiguration(), ValidationTypeMxBlacklist)
		validationDomainListMatch, result := new(validationDomainListMatchMock), validator.result
		validationRegex, validationMx, validationMxBlacklist, validationSmtp := new(validationRegexMock), new(validationMxMock), new(validationMxBlacklistMock), new(validationSmtpMock)
		validator.domainListMatch, validator.regex, validator.mx, validator.mxBlacklist, validator.smtp = validationDomainListMatch, validationRegex, validationMx, validationMxBlacklist, validationSmtp
		doPassedFromDomainListMatch(result)
		failedResult := failedValidatorResult()

		validationDomainListMatch.On("check", result).Return(result)
		validationRegex.On("check", result).Return(result)
		validationMx.On("check", result).Return(failedResult)
		validator.run()
		validationDomainListMatch.AssertExpectations(t)
		validationRegex.AssertExpectations(t)
		validationMx.AssertExpectations(t)
		validationMxBlacklist.AssertNotCalled(t, "check", failedResult)
		validationSmtp.AssertNotCalled(t, "check", failedResult)
	})

	t.Run("mx blacklist validation: domainListMatch, regex, mx succeed", func(t *testing.T) {
		validator := createValidator(randomEmail(), createConfiguration(), ValidationTypeMxBlacklist)
		validationDomainListMatch, result := new(validationDomainListMatchMock), validator.result
		validationRegex, validationMx, validationMxBlacklist, validationSmtp := new(validationRegexMock), new(validationMxMock), new(validationMxBlacklistMock), new(validationSmtpMock)
		validator.domainListMatch, validator.regex, validator.mx, validator.mxBlacklist, validator.smtp = validationDomainListMatch, validationRegex, validationMx, validationMxBlacklist, validationSmtp
		doPassedFromDomainListMatch(result)

		validationDomainListMatch.On("check", result).Return(result)
		validationRegex.On("check", result).Return(result)
		validationMx.On("check", result).Return(result)
		validationMxBlacklist.On("check", result).Return(result)
		assert.Equal(t, result, validator.run())
		validationDomainListMatch.AssertExpectations(t)
		validationRegex.AssertExpectations(t)
		validationMx.AssertExpectations(t)
		validationMxBlacklist.AssertExpectations(t)
		validationSmtp.AssertNotCalled(t, "check", result)
	})

	t.Run("smtp validation: domainListMatch succeed, regex fails", func(t *testing.T) {
		validator := createValidator(randomEmail(), createConfiguration(), ValidationTypeSMTP)
		validationDomainListMatch, result := new(validationDomainListMatchMock), validator.result
		validationRegex, validationMx, validationMxBlacklist, validationSmtp := new(validationRegexMock), new(validationMxMock), new(validationMxBlacklistMock), new(validationSmtpMock)
		validator.domainListMatch, validator.regex, validator.mx, validator.mxBlacklist, validator.smtp = validationDomainListMatch, validationRegex, validationMx, validationMxBlacklist, validationSmtp
		doPassedFromDomainListMatch(result)
		failedResult := failedValidatorResult()

		validationDomainListMatch.On("check", result).Return(result)
		validationRegex.On("check", result).Return(failedResult)
		validator.run()
		validationDomainListMatch.AssertExpectations(t)
		validationRegex.AssertExpectations(t)
		validationMx.AssertNotCalled(t, "check", failedResult)
		validationMxBlacklist.AssertNotCalled(t, "check", failedResult)
		validationSmtp.AssertNotCalled(t, "check", failedResult)
	})

	t.Run("smtp validation: domainListMatch, regex succeed, mx fails", func(t *testing.T) {
		validator := createValidator(randomEmail(), createConfiguration(), ValidationTypeSMTP)
		validationDomainListMatch, result := new(validationDomainListMatchMock), validator.result
		validationRegex, validationMx, validationMxBlacklist, validationSmtp := new(validationRegexMock), new(validationMxMock), new(validationMxBlacklistMock), new(validationSmtpMock)
		validator.domainListMatch, validator.regex, validator.mx, validator.mxBlacklist, validator.smtp = validationDomainListMatch, validationRegex, validationMx, validationMxBlacklist, validationSmtp
		doPassedFromDomainListMatch(result)
		failedResult := failedValidatorResult()

		validationDomainListMatch.On("check", result).Return(result)
		validationRegex.On("check", result).Return(result)
		validationMx.On("check", result).Return(failedResult)
		validator.run()
		validationDomainListMatch.AssertExpectations(t)
		validationRegex.AssertExpectations(t)
		validationMx.AssertExpectations(t)
		validationMxBlacklist.AssertNotCalled(t, "check", failedResult)
		validationSmtp.AssertNotCalled(t, "check", failedResult)
	})

	t.Run("smtp validation: domainListMatch, regex, mx succeed, mx blacklist fails", func(t *testing.T) {
		validator := createValidator(randomEmail(), createConfiguration(), ValidationTypeSMTP)
		validationDomainListMatch, result := new(validationDomainListMatchMock), validator.result
		validationRegex, validationMx, validationMxBlacklist, validationSmtp := new(validationRegexMock), new(validationMxMock), new(validationMxBlacklistMock), new(validationSmtpMock)
		validator.domainListMatch, validator.regex, validator.mx, validator.mxBlacklist, validator.smtp = validationDomainListMatch, validationRegex, validationMx, validationMxBlacklist, validationSmtp
		doPassedFromDomainListMatch(result)
		failedResult := failedValidatorResult()

		validationDomainListMatch.On("check", result).Return(result)
		validationRegex.On("check", result).Return(result)
		validationMx.On("check", result).Return(result)
		validationMxBlacklist.On("check", result).Return(failedResult)
		validator.run()
		validationDomainListMatch.AssertExpectations(t)
		validationRegex.AssertExpectations(t)
		validationMx.AssertExpectations(t)
		validationMxBlacklist.AssertExpectations(t)
		validationSmtp.AssertNotCalled(t, "check", result)
	})

	t.Run("smtp validation: domainListMatch, regex, mx, mx blacklist succeed", func(t *testing.T) {
		validator := createValidator(randomEmail(), createConfiguration(), ValidationTypeSMTP)
		validationDomainListMatch, result := new(validationDomainListMatchMock), validator.result
		validationRegex, validationMx, validationMxBlacklist, validationSmtp := new(validationRegexMock), new(validationMxMock), new(validationMxBlacklistMock), new(validationSmtpMock)
		validator.domainListMatch, validator.regex, validator.mx, validator.mxBlacklist, validator.smtp = validationDomainListMatch, validationRegex, validationMx, validationMxBlacklist, validationSmtp
		doPassedFromDomainListMatch(result)

		validationDomainListMatch.On("check", result).Return(result)
		validationRegex.On("check", result).Return(result)
		validationMx.On("check", result).Return(result)
		validationMxBlacklist.On("check", result).Return(result)
		validationSmtp.On("check", result).Return(result)
		assert.Equal(t, result, validator.run())
		validationDomainListMatch.AssertExpectations(t)
		validationRegex.AssertExpectations(t)
		validationMx.AssertExpectations(t)
		validationMxBlacklist.AssertExpectations(t)
		validationSmtp.AssertExpectations(t)
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

func TestAddError(t *testing.T) {
	t.Run("addes error to ValidatorResult", func(t *testing.T) {
		key, value := "some_error_key", "some_error_value"
		validatorResult := addError(new(validatorResult), key, value)

		assert.Equal(t, value, validatorResult.Errors[key])
	})
}
