package truemail

func validateSMTP(validatorResult *validatorResult) *validatorResult {
	if !validateMx(validatorResult).Success {
		return validatorResult
	}

	validatorResult.validator.addUsedValidationType(ValidationTypeSMTP)
	return validatorResult
}
