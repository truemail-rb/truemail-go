package truemail

func validateMx(validatorResult *validatorResult) *validatorResult {
	if !validateRegex(validatorResult).Success {
		return validatorResult
	}

	validatorResult.validator.addUsedValidationType(ValidationTypeMx)
	return validatorResult
}
