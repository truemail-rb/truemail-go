package truemail

// Whitelist/Blacklist validation, zero validation level
// interface implementation
func (validation *validationDomainListMatch) check(validatorResult *validatorResult) *validatorResult {
	// Failure scenario
	if validation.isBlacklistedDomain(validatorResult) ||
		(validation.isWhitelistValidation(validatorResult) && !validation.isWhitelistedDomain(validatorResult)) {
		validatorResult.ValidationType = DomainListMatchBlacklist
		validatorResult.addError(ValidationTypeDomainListMatch, DomainListMatchErrorContext)
		return validatorResult
	}

	// Successful scenario
	validatorResult.Success = true

	// Handle flow with ValidationType persisting
	if !validation.isWhitelistValidation(validatorResult) &&
		!(!validation.isBlacklistedDomain(validatorResult) && !validation.isWhitelistedDomain(validatorResult)) {
		validatorResult.ValidationType = DomainListMatchWhitelist
	}

	// Handle flow for processing validatorResult via next validation level
	if (validation.isWhitelistValidation(validatorResult) && validation.isWhitelistedDomain(validatorResult)) ||
		(!validation.isBlacklistedDomain(validatorResult) && !validation.isWhitelistedDomain(validatorResult)) {
		validatorResult.isPassFromDomainListMatch = true
	}

	return validatorResult
}

// validationDomainListMatch methods

func (validation *validationDomainListMatch) emailDomain(email string) string {
	regex, _ := newRegex(RegexDomainFromEmail)
	domainCaptureGroup := 1
	return regex.FindStringSubmatch(email)[domainCaptureGroup]
}

func (validation *validationDomainListMatch) isWhitelistedDomain(validatorResult *validatorResult) bool {
	return isIncluded(
		validatorResult.Configuration.WhitelistedDomains,
		validation.emailDomain(validatorResult.Email),
	)
}

func (validation *validationDomainListMatch) isWhitelistValidation(validatorResult *validatorResult) bool {
	return validatorResult.Configuration.WhitelistValidation
}

func (validation *validationDomainListMatch) isBlacklistedDomain(validatorResult *validatorResult) bool {
	return isIncluded(
		validatorResult.Configuration.BlacklistedDomains,
		validation.emailDomain(validatorResult.Email),
	)
}
