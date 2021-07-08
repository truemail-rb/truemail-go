package truemail

func (validation *validationMxBlacklist) check(validatorResult *validatorResult) *validatorResult {
	if isIntersected(validatorResult.Configuration.BlacklistedMxIpAddresses, validatorResult.MailServers) {
		validatorResult.Success = false
		validatorResult.addError(ValidationTypeMxBlacklist, MxBlacklistErrorContext)
	}

	return validatorResult
}

const (
	MxBlacklistErrorContext = "blacklisted mx server ip address"
)
