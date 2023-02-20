package truemail

// Regex validation, first validation level
type validationRegex struct{}

// interface implementation
func (validation *validationRegex) check(validatorResult *ValidatorResult) *ValidatorResult {
	if !validatorResult.Configuration.EmailPattern.MatchString(validatorResult.Email) {
		validatorResult.Success = false
		validatorResult.addError(validationTypeRegex, regexErrorContext)
	}

	return validatorResult
}
