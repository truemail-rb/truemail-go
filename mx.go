package truemail

func (validation *validation) mx(validatorResult *validatorResult) *validatorResult {
	if !validation.regex(validatorResult).Success {
		return validatorResult
	}

	validatorResult.validator.addUsedValidationType(ValidationTypeMx)
	return validatorResult
}
