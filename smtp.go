package truemail

func validateSMTP(validatorResult *validatorResult) *validatorResult {
	if !validateMxBlacklist(validatorResult).Success {
		return validatorResult
	}

	validatorResult.validator.addUsedValidationType(ValidationTypeSMTP)
	return validatorResult
}
