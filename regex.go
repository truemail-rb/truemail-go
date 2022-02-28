package truemail

// Regex validation, first validation level
type validationRegex struct{}

// interface implementation
func (validation *validationRegex) check(validatorResult *validatorResult) *validatorResult {
	if !validatorResult.Configuration.EmailPattern.MatchString(validatorResult.Email) {
		validatorResult.Success = false
		validatorResult.addError(validationTypeRegex, regexErrorContext)
	}

	return validatorResult
}
