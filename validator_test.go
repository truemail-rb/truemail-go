package truemail

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewValidator(t *testing.T) {
	t.Run("creates validator", func(t *testing.T) {
		email, validationType, configuration := randomEmail(), randomValidationType(), createConfiguration()
		validator := newValidator(email, validationType, configuration)
		validatorResult := validator.result

		assert.Equal(t, email, validatorResult.Email)
		assert.Equal(t, validationType, validatorResult.ValidationType)
		assert.EqualValues(t, configuration, validatorResult.Configuration)
		assert.NotSame(t, configuration, validatorResult.Configuration)
		assert.False(t, validatorResult.isPassFromDomainListMatch)
		assert.Empty(t, validatorResult.usedValidations)
	})
}

func TestValidatorValidateDomainListMatch(t *testing.T) {
	t.Run("validator#validateDomainListMatch", func(t *testing.T) {
		validator := createValidator(randomEmail(), createConfiguration())
		validationDomainListMatch, result := new(validationDomainListMatchMock), validator.result
		validator.domainListMatchLayer = validationDomainListMatch

		validationDomainListMatch.On("check", result).Return(result)
		validator.validateDomainListMatch()
		validationDomainListMatch.AssertExpectations(t)
	})
}

func TestValidatorValidateRegex(t *testing.T) {
	t.Run("validator#validateRegex", func(t *testing.T) {
		validator := createValidator(randomEmail(), createConfiguration())
		validationRegex, result := new(validationRegexMock), validator.result
		validator.regexLayer = validationRegex

		validationRegex.On("check", result).Return(result)
		validator.validateRegex()
		validationRegex.AssertExpectations(t)
		assert.Equal(t, usedValidationsByType(validationTypeRegex), validator.result.usedValidations)

	})
}
func TestValidatorValidateMx(t *testing.T) {
	t.Run("when all layers passed", func(t *testing.T) {
		validator := createValidator(randomEmail(), createConfiguration())
		validationRegex, validationMx, result := new(validationRegexMock), new(validationMxMock), validator.result
		validator.regexLayer, validator.mxLayer = validationRegex, validationMx
		result.Success = true

		validationRegex.On("check", result).Return(result)
		validationMx.On("check", result).Return(result)
		validator.validateMx()
		validationRegex.AssertExpectations(t)
		validationMx.AssertExpectations(t)
		assert.Equal(t, usedValidationsByType(validationTypeMx), validator.result.usedValidations)
	})

	t.Run("when regex layer fails", func(t *testing.T) {
		validator := createValidator(randomEmail(), createConfiguration())
		validationRegex, validationMx, result := new(validationRegexMock), new(validationMxMock), validator.result
		validator.regexLayer, validator.mxLayer = validationRegex, validationMx

		validationRegex.On("check", result).Return(result)
		validator.validateMx()
		validationRegex.AssertExpectations(t)
		validationMx.AssertNotCalled(t, "check", result)
		assert.Equal(t, usedValidationsByType(validationTypeRegex), validator.result.usedValidations)
	})
}

func TestValidatorValidateMxBlacklist(t *testing.T) {
	t.Run("when all layers passed", func(t *testing.T) {
		validator := createValidator(randomEmail(), createConfiguration())
		validationRegex, validationMx, validationMxBlacklist, result := new(validationRegexMock), new(validationMxMock), new(validationMxBlacklistMock), validator.result
		validator.regexLayer, validator.mxLayer, validator.mxBlacklistLayer = validationRegex, validationMx, validationMxBlacklist
		result.Success = true

		validationRegex.On("check", result).Return(result)
		validationMx.On("check", result).Return(result)
		validationMxBlacklist.On("check", result).Return(result)
		validator.validateMxBlacklist()
		validationRegex.AssertExpectations(t)
		validationMx.AssertExpectations(t)
		validationMxBlacklist.AssertExpectations(t)
		assert.Equal(t, usedValidationsByType(validationTypeMxBlacklist), validator.result.usedValidations)
	})

	t.Run("when regex layer fails", func(t *testing.T) {
		validator := createValidator(randomEmail(), createConfiguration())
		validationRegex, validationMx, validationMxBlacklist, result := new(validationRegexMock), new(validationMxMock), new(validationMxBlacklistMock), validator.result
		validator.regexLayer, validator.mxLayer, validator.mxBlacklistLayer = validationRegex, validationMx, validationMxBlacklist

		validationRegex.On("check", result).Return(result)
		validator.validateMxBlacklist()
		validationRegex.AssertExpectations(t)
		validationMx.AssertNotCalled(t, "check", result)
		validationMxBlacklist.AssertNotCalled(t, "check", result)
		assert.Equal(t, usedValidationsByType(validationTypeRegex), validator.result.usedValidations)
	})

	t.Run("when mx layer fails", func(t *testing.T) {
		validator := createValidator(randomEmail(), createConfiguration())
		validationRegex, validationMx, validationMxBlacklist, result := new(validationRegexMock), new(validationMxMock), new(validationMxBlacklistMock), validator.result
		validator.regexLayer, validator.mxLayer, validator.mxBlacklistLayer = validationRegex, validationMx, validationMxBlacklist
		result.Success = true
		failedResult := failedValidatorResult()

		validationRegex.On("check", result).Return(result)
		validationMx.On("check", result).Return(failedResult)
		validator.validateMxBlacklist()
		validationRegex.AssertExpectations(t)
		validationMx.AssertExpectations(t)
		validationMxBlacklist.AssertNotCalled(t, "check", failedResult)
		assert.Equal(t, usedValidationsByType(validationTypeMx), validator.result.usedValidations)
	})
}

func TestValidatorValidateSMTP(t *testing.T) {
	t.Run("when all layers passed", func(t *testing.T) {
		validator := createValidator(randomEmail(), createConfiguration())
		validationRegex, validationMx, validationMxBlacklist, validationSmtp := new(validationRegexMock), new(validationMxMock), new(validationMxBlacklistMock), new(validationSmtpMock)
		validator.regexLayer, validator.mxLayer, validator.mxBlacklistLayer, validator.smtpLayer = validationRegex, validationMx, validationMxBlacklist, validationSmtp
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
		assert.Equal(t, usedValidationsByType(validationTypeSmtp), validator.result.usedValidations)
	})

	t.Run("when regex layer fails", func(t *testing.T) {
		validator := createValidator(randomEmail(), createConfiguration())
		validationRegex, validationMx, validationMxBlacklist, validationSmtp := new(validationRegexMock), new(validationMxMock), new(validationMxBlacklistMock), new(validationSmtpMock)
		validator.regexLayer, validator.mxLayer, validator.mxBlacklistLayer, validator.smtpLayer = validationRegex, validationMx, validationMxBlacklist, validationSmtp
		result := validator.result
		result.Success = true
		failedResult := failedValidatorResult()

		validationRegex.On("check", result).Return(failedResult)
		validator.validateSMTP()
		validationRegex.AssertExpectations(t)
		validationMx.AssertNotCalled(t, "check", failedResult)
		validationMxBlacklist.AssertNotCalled(t, "check", failedResult)
		validationSmtp.AssertNotCalled(t, "check", failedResult)
		assert.Equal(t, usedValidationsByType(validationTypeRegex), validator.result.usedValidations)
	})

	t.Run("when mx layer fails", func(t *testing.T) {
		validator := createValidator(randomEmail(), createConfiguration())
		validationRegex, validationMx, validationMxBlacklist, validationSmtp := new(validationRegexMock), new(validationMxMock), new(validationMxBlacklistMock), new(validationSmtpMock)
		validator.regexLayer, validator.mxLayer, validator.mxBlacklistLayer, validator.smtpLayer = validationRegex, validationMx, validationMxBlacklist, validationSmtp
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
		assert.Equal(t, usedValidationsByType(validationTypeMx), validator.result.usedValidations)
	})

	t.Run("when mx blacklist layer fails", func(t *testing.T) {
		validator := createValidator(randomEmail(), createConfiguration())
		validationRegex, validationMx, validationMxBlacklist, validationSmtp := new(validationRegexMock), new(validationMxMock), new(validationMxBlacklistMock), new(validationSmtpMock)
		validator.regexLayer, validator.mxLayer, validator.mxBlacklistLayer, validator.smtpLayer = validationRegex, validationMx, validationMxBlacklist, validationSmtp
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
		assert.Equal(t, usedValidationsByType(validationTypeMxBlacklist), validator.result.usedValidations)
	})

	t.Run("when smtp layer fails", func(t *testing.T) {
		validator := createValidator(randomEmail(), createConfiguration())
		validationRegex, validationMx, validationMxBlacklist, validationSmtp := new(validationRegexMock), new(validationMxMock), new(validationMxBlacklistMock), new(validationSmtpMock)
		validator.regexLayer, validator.mxLayer, validator.mxBlacklistLayer, validator.smtpLayer = validationRegex, validationMx, validationMxBlacklist, validationSmtp
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
		assert.Equal(t, usedValidationsByType(validationTypeSmtp), validator.result.usedValidations)
	})
}

func TestValidatorRun(t *testing.T) {
	t.Run("domainListMatchLayer fails", func(t *testing.T) {
		validator := createValidator(randomEmail(), createConfiguration(), validationTypeRegex)
		validationDomainListMatch, result := new(validationDomainListMatchMock), validator.result
		validationRegex, validationMx, validationMxBlacklist, validationSmtp := new(validationRegexMock), new(validationMxMock), new(validationMxBlacklistMock), new(validationSmtpMock)
		validator.domainListMatchLayer, validator.regexLayer, validator.mxLayer, validator.mxBlacklistLayer, validator.smtpLayer = validationDomainListMatch, validationRegex, validationMx, validationMxBlacklist, validationSmtp

		validationDomainListMatch.On("check", result).Return(result)
		assert.Equal(t, result, validator.run())
		validationDomainListMatch.AssertExpectations(t)
		validationRegex.AssertNotCalled(t, "check", result)
		validationMx.AssertNotCalled(t, "check", result)
		validationMxBlacklist.AssertNotCalled(t, "check", result)
		validationSmtp.AssertNotCalled(t, "check", result)
	})

	t.Run("regex validation: domainListMatchLayer succeed", func(t *testing.T) {
		validator := createValidator(randomEmail(), createConfiguration(), validationTypeRegex)
		validationDomainListMatch, result := new(validationDomainListMatchMock), validator.result
		validationRegex, validationMx, validationMxBlacklist, validationSmtp := new(validationRegexMock), new(validationMxMock), new(validationMxBlacklistMock), new(validationSmtpMock)
		validator.domainListMatchLayer, validator.regexLayer, validator.mxLayer, validator.mxBlacklistLayer, validator.smtpLayer = validationDomainListMatch, validationRegex, validationMx, validationMxBlacklist, validationSmtp
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

	t.Run("mx validation: domainListMatchLayer succeed, regex fails", func(t *testing.T) {
		validator := createValidator(randomEmail(), createConfiguration(), validationTypeMx)
		validationDomainListMatch, result := new(validationDomainListMatchMock), validator.result
		validationRegex, validationMx, validationMxBlacklist, validationSmtp := new(validationRegexMock), new(validationMxMock), new(validationMxBlacklistMock), new(validationSmtpMock)
		validator.domainListMatchLayer, validator.regexLayer, validator.mxLayer, validator.mxBlacklistLayer, validator.smtpLayer = validationDomainListMatch, validationRegex, validationMx, validationMxBlacklist, validationSmtp
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

	t.Run("mx validation: domainListMatchLayer, regex succeed", func(t *testing.T) {
		validator := createValidator(randomEmail(), createConfiguration(), validationTypeMx)
		validationDomainListMatch, result := new(validationDomainListMatchMock), validator.result
		validationRegex, validationMx, validationMxBlacklist, validationSmtp := new(validationRegexMock), new(validationMxMock), new(validationMxBlacklistMock), new(validationSmtpMock)
		validator.domainListMatchLayer, validator.regexLayer, validator.mxLayer, validator.mxBlacklistLayer, validator.smtpLayer = validationDomainListMatch, validationRegex, validationMx, validationMxBlacklist, validationSmtp
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

	t.Run("mx blacklist validation: domainListMatchLayer succeed, regex fails", func(t *testing.T) {
		validator := createValidator(randomEmail(), createConfiguration(), validationTypeMxBlacklist)
		validationDomainListMatch, result := new(validationDomainListMatchMock), validator.result
		validationRegex, validationMx, validationMxBlacklist, validationSmtp := new(validationRegexMock), new(validationMxMock), new(validationMxBlacklistMock), new(validationSmtpMock)
		validator.domainListMatchLayer, validator.regexLayer, validator.mxLayer, validator.mxBlacklistLayer, validator.smtpLayer = validationDomainListMatch, validationRegex, validationMx, validationMxBlacklist, validationSmtp
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

	t.Run("mx blacklist validation: domainListMatchLayer, regex succeed, mx fails", func(t *testing.T) {
		validator := createValidator(randomEmail(), createConfiguration(), validationTypeMxBlacklist)
		validationDomainListMatch, result := new(validationDomainListMatchMock), validator.result
		validationRegex, validationMx, validationMxBlacklist, validationSmtp := new(validationRegexMock), new(validationMxMock), new(validationMxBlacklistMock), new(validationSmtpMock)
		validator.domainListMatchLayer, validator.regexLayer, validator.mxLayer, validator.mxBlacklistLayer, validator.smtpLayer = validationDomainListMatch, validationRegex, validationMx, validationMxBlacklist, validationSmtp
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

	t.Run("mx blacklist validation: domainListMatchLayer, regex, mx succeed", func(t *testing.T) {
		validator := createValidator(randomEmail(), createConfiguration(), validationTypeMxBlacklist)
		validationDomainListMatch, result := new(validationDomainListMatchMock), validator.result
		validationRegex, validationMx, validationMxBlacklist, validationSmtp := new(validationRegexMock), new(validationMxMock), new(validationMxBlacklistMock), new(validationSmtpMock)
		validator.domainListMatchLayer, validator.regexLayer, validator.mxLayer, validator.mxBlacklistLayer, validator.smtpLayer = validationDomainListMatch, validationRegex, validationMx, validationMxBlacklist, validationSmtp
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

	t.Run("smtp validation: domainListMatchLayer succeed, regex fails", func(t *testing.T) {
		validator := createValidator(randomEmail(), createConfiguration(), validationTypeSmtp)
		validationDomainListMatch, result := new(validationDomainListMatchMock), validator.result
		validationRegex, validationMx, validationMxBlacklist, validationSmtp := new(validationRegexMock), new(validationMxMock), new(validationMxBlacklistMock), new(validationSmtpMock)
		validator.domainListMatchLayer, validator.regexLayer, validator.mxLayer, validator.mxBlacklistLayer, validator.smtpLayer = validationDomainListMatch, validationRegex, validationMx, validationMxBlacklist, validationSmtp
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

	t.Run("smtp validation: domainListMatchLayer, regex succeed, mx fails", func(t *testing.T) {
		validator := createValidator(randomEmail(), createConfiguration(), validationTypeSmtp)
		validationDomainListMatch, result := new(validationDomainListMatchMock), validator.result
		validationRegex, validationMx, validationMxBlacklist, validationSmtp := new(validationRegexMock), new(validationMxMock), new(validationMxBlacklistMock), new(validationSmtpMock)
		validator.domainListMatchLayer, validator.regexLayer, validator.mxLayer, validator.mxBlacklistLayer, validator.smtpLayer = validationDomainListMatch, validationRegex, validationMx, validationMxBlacklist, validationSmtp
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

	t.Run("smtp validation: domainListMatchLayer, regex, mx succeed, mx blacklist fails", func(t *testing.T) {
		validator := createValidator(randomEmail(), createConfiguration(), validationTypeSmtp)
		validationDomainListMatch, result := new(validationDomainListMatchMock), validator.result
		validationRegex, validationMx, validationMxBlacklist, validationSmtp := new(validationRegexMock), new(validationMxMock), new(validationMxBlacklistMock), new(validationSmtpMock)
		validator.domainListMatchLayer, validator.regexLayer, validator.mxLayer, validator.mxBlacklistLayer, validator.smtpLayer = validationDomainListMatch, validationRegex, validationMx, validationMxBlacklist, validationSmtp
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

	t.Run("smtp validation: domainListMatchLayer, regex, mx, mx blacklist succeed", func(t *testing.T) {
		validator := createValidator(randomEmail(), createConfiguration(), validationTypeSmtp)
		validationDomainListMatch, result := new(validationDomainListMatchMock), validator.result
		validationRegex, validationMx, validationMxBlacklist, validationSmtp := new(validationRegexMock), new(validationMxMock), new(validationMxBlacklistMock), new(validationSmtpMock)
		validator.domainListMatchLayer, validator.regexLayer, validator.mxLayer, validator.mxBlacklistLayer, validator.smtpLayer = validationDomainListMatch, validationRegex, validationMx, validationMxBlacklist, validationSmtp
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

func TestValidatorResultAddUsedValidationType(t *testing.T) {
	t.Run("validatorResult#addUsedValidationType", func(t *testing.T) {
		validationType := randomValidationType()
		result := new(ValidatorResult)
		result.addUsedValidationType(validationType)

		assert.Equal(t, []string{validationType}, result.usedValidations)
	})
}

func TestValidatorResultAddError(t *testing.T) {
	t.Run("validatorResult#addError", func(t *testing.T) {
		key, value := "some_error_key", "some_error_value"
		result := new(ValidatorResult)
		result.addError(key, value)

		assert.Equal(t, value, result.Errors[key])
	})
}
