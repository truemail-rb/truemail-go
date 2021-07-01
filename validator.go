package truemail

import "fmt"

func Validate(email string, configuration *configuration, options ...string) (*validatorResult, error) {
	validationType, err := variadicValidationType(options)

	if err != nil {
		return nil, err
	}

	return newValidator(email, validationType, configuration).run(), err
}

type validationDomainListMatch struct{}
type validationRegex struct{}
type validationMx struct{}
type validationMxBlacklist struct{}
type validationSmtp struct{}

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

// validatorResult structure
type validatorResult struct {
	Success, SMTPDebug, isPassFromDomainListMatch bool
	Email, Domain, ValidationType                 string
	MailServers, usedValidations                  []string
	Errors                                        map[string]string
	Configuration                                 *configuration
}

// validator, structure with behaviour
type validator struct {
	result *validatorResult
	domainListMatch
	regex
	mx
	mxBlacklist
	smtp
}

func variadicValidationType(options []string) (string, error) {
	if len(options) == 0 {
		return ValidationTypeDefault, nil
	}

	validationType := options[0]
	return validationType, validateValidationTypeContext(validationType)
}

func validateValidationTypeContext(validationType string) error {
	if isIncluded(availableValidationTypes(), validationType) {
		return nil
	}
	return fmt.Errorf(
		"%s is invalid validation type, use one of these: %s",
		validationType,
		availableValidationTypes(),
	)
}

func newValidator(email, validationType string, configuration *configuration) *validator {
	validator := &validator{
		result: &validatorResult{
			Email:          email,
			Configuration:  configuration,
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

func addError(validatorResult *validatorResult, key, value string) *validatorResult {
	if validatorResult.Errors == nil {
		validatorResult.Errors = map[string]string{}
	}
	validatorResult.Errors[key] = value
	return validatorResult
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

// validatorResult methods

func (validatorResult *validatorResult) addUsedValidationType(validationType string) {
	validatorResult.usedValidations = append(validatorResult.usedValidations, validationType)
}
