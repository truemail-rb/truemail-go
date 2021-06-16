package truemail

// Whitelist/Blacklist validation, zero validation level
func validateDomainListMatch(validatorResult *validatorResult) *validatorResult {
	// Failure scenario
	if isBlacklistedDomain(validatorResult) ||
		(isWhitelistValidation(validatorResult) && !isWhitelistedDomain(validatorResult)) {
		validatorResult.ValidationType = DomainListMatchBlacklist
		return addError(validatorResult, DomainListMatchErrorKey, DomainListMatchErrorContext)
	}

	// Successful scenario
	validatorResult.Success = true

	// Handle flow with ValidationType persisting
	if !isWhitelistValidation(validatorResult) &&
		!(!isBlacklistedDomain(validatorResult) && !isWhitelistedDomain(validatorResult)) {
		validatorResult.ValidationType = DomainListMatchWhitelist
	}

	// Handle flow for processing validatorResult via next validation level
	if (isWhitelistValidation(validatorResult) && isWhitelistedDomain(validatorResult)) ||
		(!isBlacklistedDomain(validatorResult) && !isWhitelistedDomain(validatorResult)) {
		validatorResult.validator.isPassFromDomainListMatch = true
	}

	return validatorResult
}

const (
	DomainListMatchWhitelist    = "whitelist"
	DomainListMatchBlacklist    = "blacklist"
	DomainListMatchErrorKey     = "domain_list_match"
	DomainListMatchErrorContext = "blacklisted email"
)

func emailDomain(email string) string {
	regex, _ := newRegex(RegexDomainFromEmail)
	domainCaptureGroup := 1
	return regex.FindStringSubmatch(email)[domainCaptureGroup]
}

func isWhitelistedDomain(validatorResult *validatorResult) bool {
	return isIncluded(
		validatorResult.Configuration.WhitelistedDomains,
		emailDomain(validatorResult.Email),
	)
}

func isWhitelistValidation(validatorResult *validatorResult) bool {
	return validatorResult.Configuration.WhitelistValidation
}

func isBlacklistedDomain(validatorResult *validatorResult) bool {
	return isIncluded(
		validatorResult.Configuration.BlacklistedDomains,
		emailDomain(validatorResult.Email),
	)
}
