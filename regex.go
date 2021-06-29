package truemail

func (validation *validation) regex(validatorResult *validatorResult) *validatorResult {
	validatorResult.validator.addUsedValidationType(ValidationTypeRegex)
	return validatorResult
}
