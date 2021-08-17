package truemail

// SMTP validation, fourth validation level
type validationSmtp struct{}

// interface implementation
func (validation *validationSmtp) check(validatorResult *validatorResult) *validatorResult {
	return validatorResult
}
