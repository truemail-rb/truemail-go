package truemail

func (validation *validation) mxBlacklist(validatorResult *validatorResult) *validatorResult {
	if !validation.mx(validatorResult).Success {
		return validatorResult
	}

	return validatorResult
}
