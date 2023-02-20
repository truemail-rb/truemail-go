package truemail

// SMTP validation, fourth validation level
type validationSmtp struct {
	result      *ValidatorResult
	smtpResults []*SmtpRequest
	builder
}

// interface implementation
func (validation *validationSmtp) check(validatorResult *ValidatorResult) *ValidatorResult {
	validation.result = validatorResult
	validation.initSmtpBuilder()
	validation.run()

	if validation.isIncludesSuccessfulSmtpResponse() {
		return validatorResult
	}

	validation.result.SmtpDebug = validation.smtpResults

	if validation.isSmtpSafeCheckEnabled() && validation.isNotIncludeUserNotFoundErrors() {
		return validatorResult
	}

	validatorResult.Success = false
	validatorResult.addError(validationTypeSmtp, smtpErrorContext)

	return validatorResult
}

// validationSmtp methods

// Initializes SMTP validation SMTP entities builder
func (validation *validationSmtp) initSmtpBuilder() {
	validation.builder = new(smtpBuilder)
}

// Runs SMTP session for each target server until receive successful session response
func (validation *validationSmtp) run() {
	for _, targetHostAddress := range validation.filteredMailServersByFailFastScenario() {
		if validation.runSmtpSession(targetHostAddress) {
			break
		}
	}
}

// Runs SMTP session for target mail server. Returns
// true for successful session, otherwise returns false
func (validation *validationSmtp) runSmtpSession(targetHostAddress string) bool {
	validatorResult, validatorBuilder := validation.result, validation.builder
	smtpRequest := validatorBuilder.newSmtpRequest(
		validation.attempts(),
		validatorResult.Email,
		targetHostAddress,
		validatorResult.Configuration,
	)
	smtpResponse := smtpRequest.Response
	validation.smtpResults = append(validation.smtpResults, smtpRequest)

	for smtpRequest.Attempts > 0 {
		smtpClient := validatorBuilder.newSmtpClient(smtpRequest.Configuration)
		smtpRequest.Attempts -= 1

		if smtpClient.runSession() {
			smtpResponse.Rcptto = true
			return true
		}

		smtpResponse.Errors = append(smtpResponse.Errors, smtpClient.sessionError())
	}

	return false
}

// Returns true if SMTP fail fast scenario is enabled, otherwise returns false
func (validation *validationSmtp) isFailFastScenario() bool {
	return validation.result.Configuration.SmtpFailFast
}

// Returns first item from validationResult.MailServers if SMTP fail fast scenario is enabled,
// otherwise returns all validationResult.MailServers items
func (validation *validationSmtp) filteredMailServersByFailFastScenario() []string {
	mailServers := validation.result.MailServers

	if validation.isFailFastScenario() {
		return mailServers[:1]
	}

	return mailServers
}

// Returns true for case when more than one mail server exists, otherwise returns false
func (validation *validationSmtp) isMoreThanOneMailServer() bool {
	return len(validation.result.MailServers) > 1
}

// Returns 1 for SMTP fail fast scenario or for case when more than one mail server exists,
// otherwise returns number of connection attempts defined in configuration
func (validation *validationSmtp) attempts() int {
	if validation.isFailFastScenario() || validation.isMoreThanOneMailServer() {
		return 1
	}

	return validation.result.Configuration.ConnectionAttempts
}

// Terminates iteration and returns true for empty slice or when first successful SMTP response found,
// returns false if successful SMTP response not found
func (validation *validationSmtp) isIncludesSuccessfulSmtpResponse() (successfulSmtpResponse bool) {
	smtpResults := validation.smtpResults

	if len(smtpResults) == 0 {
		return true
	}

	for _, smtpRequest := range smtpResults {
		if !smtpRequest.Response.Rcptto {
			continue
		}
		successfulSmtpResponse = true
		break
	}

	return successfulSmtpResponse
}

// Returns true if SMTP safe check scenario is enabled, otherwise returns false
func (validation *validationSmtp) isSmtpSafeCheckEnabled() bool {
	return validation.result.Configuration.SmtpSafeCheck
}

// Returns true if SMTP results does not contain UserNotFound erros,
// otherwise terminates iteration and returns false
func (validation *validationSmtp) isNotIncludeUserNotFoundErrors() bool {
	for _, smtpRequest := range validation.smtpResults {
		for _, err := range smtpRequest.Response.Errors {
			if err.isRecptTo && validation.result.Configuration.SmtpErrorBodyPattern.MatchString(err.Error()) {
				return false
			}
		}
	}
	return true
}
