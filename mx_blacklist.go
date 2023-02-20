package truemail

// MX blacklist validation, third validation level
type validationMxBlacklist struct{}

// interface implementation
func (validation *validationMxBlacklist) check(validatorResult *ValidatorResult) *ValidatorResult {
	if isIntersected(validatorResult.Configuration.BlacklistedMxIpAddresses, validatorResult.MailServers) {
		validatorResult.Success = false
		validatorResult.addError(validationTypeMxBlacklist, mxBlacklistErrorContext)
	}

	return validatorResult
}
