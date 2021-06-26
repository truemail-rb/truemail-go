package truemail

func validateMxBlacklist(validatorResult *validatorResult) *validatorResult {
	if !validateMx(validatorResult).Success {
		return validatorResult
	}

	return validatorResult
}
