package truemail

func (validation *validationRegex) check(validatorResult *validatorResult) *validatorResult {
	if !validatorResult.Configuration.EmailPattern.MatchString(validatorResult.Email) {
		validatorResult.Success = false
		validatorResult.addError(ValidationTypeRegex, RegexErrorContext)
	}

	return validatorResult
}
