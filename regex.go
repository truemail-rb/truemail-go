package truemail

func validateRegex(validatorResult *validatorResult) *validatorResult {
	validatorResult.validator.addUsedValidationType(ValidationTypeRegex)
	return validatorResult
}
