package truemail

func (validation *validationRegex) check(validatorResult *validatorResult) *validatorResult {
	result := validatorResult.Configuration.EmailPattern.MatchString(validatorResult.Email)
	validatorResult.Success = result

	if !result {
		validatorResult.addError(ValidationTypeRegex, RegexErrorContext)
	}

	return validatorResult
}

const (
	RegexErrorContext = "email does not match the regular expression"
)
