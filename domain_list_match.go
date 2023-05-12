package truemail

// Whitelist/Blacklist validation, zero validation level
type validationDomainListMatch struct{ result *ValidatorResult }

// interface implementation
func (validation *validationDomainListMatch) check(validatorResult *ValidatorResult) *ValidatorResult {
	validation.result = validatorResult
	validation.setValidatorResultDomain()

	// Failure scenario
	if validation.isBlacklistedDomain() || (validation.isWhitelistValidation() && !validation.isWhitelistedDomain()) {
		validatorResult.ValidationType = domainListMatchBlacklist
		validatorResult.addError(validationTypeDomainListMatch, domainListMatchErrorContext)
		return validatorResult
	}

	// Successful scenario
	validatorResult.Success = true

	// Handle flow with ValidationType persisting
	if !validation.isWhitelistValidation() && !(!validation.isBlacklistedDomain() && !validation.isWhitelistedDomain()) {
		validatorResult.ValidationType = domainListMatchWhitelist
	}

	// Handle flow for processing validatorResult via next validation level
	if (validation.isWhitelistValidation() && validation.isWhitelistedDomain()) ||
		(!validation.isBlacklistedDomain() && !validation.isWhitelistedDomain()) {
		validatorResult.isPassFromDomainListMatch = true
	}

	return validatorResult
}

// validationDomainListMatch methods

// Assigns domain based on validator result email to validatorResult
func (validation *validationDomainListMatch) setValidatorResultDomain() {
	validatorResult := validation.result
	validatorResult.Domain = emailDomain(validatorResult.Email)
}

// Returns true if email domain is included in whitelisted domains slice, otherwise returns false
func (validation *validationDomainListMatch) isWhitelistedDomain() bool {
	validatorResult := validation.result
	return isIncluded(validatorResult.Configuration.WhitelistedDomains, validatorResult.Domain)
}

// Returns true if whitelist validation enabled, otherwise returns false
func (validation *validationDomainListMatch) isWhitelistValidation() bool {
	validatorResult := validation.result
	return validatorResult.Configuration.WhitelistValidation
}

// Returns true if email domain is included in blacklisted domains slice, otherwise returns false
func (validation *validationDomainListMatch) isBlacklistedDomain() bool {
	validatorResult := validation.result
	return isIncluded(validatorResult.Configuration.BlacklistedDomains, validatorResult.Domain)
}
