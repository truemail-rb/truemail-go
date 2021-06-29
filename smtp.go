package truemail

func (validation *validation) smtp(validatorResult *validatorResult) *validatorResult {
	if !validation.mxBlacklist(validatorResult).Success {
		return validatorResult
	}

	validatorResult.validator.addUsedValidationType(ValidationTypeSMTP)
	return validatorResult
}
