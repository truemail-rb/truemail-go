package truemail

// validation layers structures

type validationDomainListMatch struct{}
type validationRegex struct{}
type validationMxBlacklist struct{}
type validationSmtp struct{}

// Validator result mutable object. Each validation
// layer write something into validatorResult
type validatorResult struct {
	Success, SMTPDebug, isPassFromDomainListMatch                bool
	Email, Domain, ValidationType, punycodeEmail, punycodeDomain string
	MailServers, usedValidations                                 []string
	Errors                                                       map[string]string
	Configuration                                                *configuration
}

// validatorResult methods

func (validatorResult *validatorResult) addUsedValidationType(validationType string) {
	validatorResult.usedValidations = append(validatorResult.usedValidations, validationType)
}

func (validatorResult *validatorResult) addError(key, value string) {
	if validatorResult.Errors == nil {
		validatorResult.Errors = map[string]string{}
	}
	validatorResult.Errors[key] = value
}

// Structure with behaviour. Responsible for the
// logic of calling the validation layers sequence
type validator struct {
	result *validatorResult
	domainListMatch
	regex
	mx
	mxBlacklist
	smtp
}

// New validator builder. Returns consistent validator structure
func newValidator(email, validationType string, configuration *configuration) *validator {
	validator := &validator{
		result: &validatorResult{
			Email:          email,
			Configuration:  copyConfigurationByPointer(configuration),
			ValidationType: validationType,
		},
		domainListMatch: &validationDomainListMatch{},
		regex:           &validationRegex{},
		mx:              &validationMx{},
		mxBlacklist:     &validationMxBlacklist{},
		smtp:            &validationSmtp{},
	}

	return validator
}

// validation layers interfaces

type domainListMatch interface {
	check(validatorResult *validatorResult) *validatorResult
}

type regex interface {
	check(validatorResult *validatorResult) *validatorResult
}

type mx interface {
	check(validatorResult *validatorResult) *validatorResult
}

type mxBlacklist interface {
	check(validatorResult *validatorResult) *validatorResult
}

type smtp interface {
	check(validatorResult *validatorResult) *validatorResult
}

// validator methods

func (validator *validator) validateDomainListMatch() {
	validator.domainListMatch.check(validator.result)
}

func (validator *validator) validateRegex() {
	validatorResult := validator.result
	validatorResult.addUsedValidationType(ValidationTypeRegex)
	validator.regex.check(validatorResult)
}

func (validator *validator) validateMx() {
	validatorResult := validator.result

	validatorResult.addUsedValidationType(ValidationTypeRegex)
	if !validator.regex.check(validatorResult).Success {
		return
	}

	validatorResult.addUsedValidationType(ValidationTypeMx)
	validator.mx.check(validatorResult)
}

func (validator *validator) validateMxBlacklist() {
	validatorResult := validator.result

	validatorResult.addUsedValidationType(ValidationTypeRegex)
	if !validator.regex.check(validatorResult).Success {
		return
	}

	validatorResult.addUsedValidationType(ValidationTypeMx)
	if !validator.mx.check(validatorResult).Success {
		return
	}

	validatorResult.addUsedValidationType(ValidationTypeMxBlacklist)
	validator.mxBlacklist.check(validatorResult)
}

func (validator *validator) validateSMTP() {
	validatorResult := validator.result

	validatorResult.addUsedValidationType(ValidationTypeRegex)
	if !validator.regex.check(validatorResult).Success {
		return
	}

	validatorResult.addUsedValidationType(ValidationTypeMx)
	if !validator.mx.check(validatorResult).Success {
		return
	}

	validatorResult.addUsedValidationType(ValidationTypeMxBlacklist)
	if !validator.mxBlacklist.check(validatorResult).Success {
		return
	}

	validatorResult.addUsedValidationType(ValidationTypeSMTP)
	validator.smtp.check(validatorResult)
}

// validator entrypoint. This method triggers chain of
// validation layers
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
	case ValidationTypeRegex:
		validator.validateRegex()
	case ValidationTypeMx:
		validator.validateMx()
	case ValidationTypeMxBlacklist:
		validator.validateMxBlacklist()
	case ValidationTypeSMTP:
		validator.validateSMTP()
	}
	return validatorResult
}
