package truemail

// Validator result mutable structure. Each validation
// layer write something into validatorResult
type validatorResult struct {
	Success, isPassFromDomainListMatch                           bool
	Email, Domain, ValidationType, punycodeEmail, punycodeDomain string
	MailServers, usedValidations                                 []string
	Errors                                                       map[string]string
	Configuration                                                *configuration
	SmtpDebug                                                    []*smtpRequest
}

// validatorResult methods

// Addes current validation type to validator result used validations slice
func (validatorResult *validatorResult) addUsedValidationType(validationType string) {
	validatorResult.usedValidations = append(validatorResult.usedValidations, validationType)
}

// Addes error to validator result errors dictionary
func (validatorResult *validatorResult) addError(key, value string) {
	if validatorResult.Errors == nil {
		validatorResult.Errors = map[string]string{}
	}
	validatorResult.Errors[key] = value
}

// Structure with behavior. Responsible for the
// logic of calling the validation layers sequence
type validator struct {
	result *validatorResult
	domainListMatchLayer
	regexLayer
	mxLayer
	mxBlacklistLayer
	smtpLayer
}

// New validator builder. Returns consistent validator structure
func newValidator(email, validationType string, configuration *configuration) *validator {
	validator := &validator{
		result: &validatorResult{
			Email:          email,
			Configuration:  copyConfigurationByPointer(configuration),
			ValidationType: validationType,
		},
		domainListMatchLayer: &validationDomainListMatch{},
		regexLayer:           &validationRegex{},
		mxLayer:              &validationMx{},
		mxBlacklistLayer:     &validationMxBlacklist{},
		smtpLayer:            &validationSmtp{},
	}

	return validator
}

// validation layers interfaces

type domainListMatchLayer interface {
	check(validatorResult *validatorResult) *validatorResult
}

type regexLayer interface {
	check(validatorResult *validatorResult) *validatorResult
}

type mxLayer interface {
	check(validatorResult *validatorResult) *validatorResult
}

type mxBlacklistLayer interface {
	check(validatorResult *validatorResult) *validatorResult
}

type smtpLayer interface {
	check(validatorResult *validatorResult) *validatorResult
}

// validator methods

// Runs Whitelist/Blacklist validation
func (validator *validator) validateDomainListMatch() {
	validator.domainListMatchLayer.check(validator.result)
}

// Runs Regex validation
func (validator *validator) validateRegex() {
	validatorResult := validator.result
	validatorResult.addUsedValidationType(validationTypeRegex)
	validator.regexLayer.check(validatorResult)
}

// Runs validations chain: Regex -> Mx
func (validator *validator) validateMx() {
	validatorResult := validator.result

	validatorResult.addUsedValidationType(validationTypeRegex)
	if !validator.regexLayer.check(validatorResult).Success {
		return
	}

	validatorResult.addUsedValidationType(validationTypeMx)
	validator.mxLayer.check(validatorResult)
}

// Runs validations chain: Regex -> Mx -> MxBlacklist
func (validator *validator) validateMxBlacklist() {
	validatorResult := validator.result

	validatorResult.addUsedValidationType(validationTypeRegex)
	if !validator.regexLayer.check(validatorResult).Success {
		return
	}

	validatorResult.addUsedValidationType(validationTypeMx)
	if !validator.mxLayer.check(validatorResult).Success {
		return
	}

	validatorResult.addUsedValidationType(validationTypeMxBlacklist)
	validator.mxBlacklistLayer.check(validatorResult)
}

// Runs validations chain: Regex -> Mx -> MxBlacklist -> SMTP
func (validator *validator) validateSMTP() {
	validatorResult := validator.result

	validatorResult.addUsedValidationType(validationTypeRegex)
	if !validator.regexLayer.check(validatorResult).Success {
		return
	}

	validatorResult.addUsedValidationType(validationTypeMx)
	if !validator.mxLayer.check(validatorResult).Success {
		return
	}

	validatorResult.addUsedValidationType(validationTypeMxBlacklist)
	if !validator.mxBlacklistLayer.check(validatorResult).Success {
		return
	}

	validatorResult.addUsedValidationType(validationTypeSmtp)
	validator.smtpLayer.check(validatorResult)
}

// validator entrypoint. This method triggers chain of validation layers
func (validator *validator) run() *validatorResult {
	// TODO: add painc if run will called more then one time
	// or check len(validatorResult.usedValidations) == 0

	// preparing for running
	validatorResult := validator.result
	validatorResult.usedValidations = []string{}

	// Whitelist/Blacklist validation
	validator.validateDomainListMatch()
	if !validatorResult.Success || !validatorResult.isPassFromDomainListMatch {
		return validatorResult
	}
	// run validation flow
	switch validatorResult.ValidationType {
	case validationTypeRegex:
		validator.validateRegex()
	case validationTypeMx:
		validator.validateMx()
	case validationTypeMxBlacklist:
		validator.validateMxBlacklist()
	case validationTypeSmtp:
		validator.validateSMTP()
	}
	return validatorResult
}
