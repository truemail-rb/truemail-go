package truemail

// MX blacklist validation, third validation level
// interface implementation
func (validation *validationMxBlacklist) check(validatorResult *validatorResult) *validatorResult {
	if isIntersected(validatorResult.Configuration.BlacklistedMxIpAddresses, validatorResult.MailServers) {
		validatorResult.Success = false
		validatorResult.addError(ValidationTypeMxBlacklist, MxBlacklistErrorContext)
	}

	return validatorResult
}
