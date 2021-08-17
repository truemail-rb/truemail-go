package truemail

// Whitelist/Blacklist validation, zero validation level
type validationDomainListMatch struct{}

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

// Returns true if email domain is included in whitelisted domains slice, otherwise returns false
func (validation *validationDomainListMatch) isWhitelistedDomain(validatorResult *validatorResult) bool {
	return isIncluded(
		validatorResult.Configuration.WhitelistedDomains,
		emailDomain(validatorResult.Email),
	)
}

// Returns true if whitelist validation enebled, otherwise returns false
func (validation *validationDomainListMatch) isWhitelistValidation(validatorResult *validatorResult) bool {
	return validatorResult.Configuration.WhitelistValidation
}

// Returns true if email domain is included in blacklisted domains slice, otherwise returns false
func (validation *validationDomainListMatch) isBlacklistedDomain(validatorResult *validatorResult) bool {
	return isIncluded(
		validatorResult.Configuration.BlacklistedDomains,
		emailDomain(validatorResult.Email),
	)
}
