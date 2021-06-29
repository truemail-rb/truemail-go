package truemail

import "fmt"

func Validate(email string, configuration *configuration, options ...string) (*validatorResult, error) {
	validationType, err := variadicValidationType(options)

	if err != nil {
		return nil, err
	}

	return newValidator(email, validationType, configuration).run(), err
}

type validate interface {
	domainListMatch(validatorResult *validatorResult) *validatorResult
	regex(validatorResult *validatorResult) *validatorResult
	mx(validatorResult *validatorResult) *validatorResult
	smtp(validatorResult *validatorResult) *validatorResult
}

// validatorResult structure
type validatorResult struct {
	Success, SMTPDebug            bool
	Email, Domain, ValidationType string
	MailServers                   []string
	Errors                        map[string]string
	Configuration                 *configuration
	validator                     *validator
}

// validation structure with bunch of methods
type validation struct{}

// validator, structure with behaviour
type validator struct {
	result                    *validatorResult
	usedValidations           []string
	isPassFromDomainListMatch bool
	validate
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
		validate: &validation{},
	}

	validator.result.validator = validator

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
	validator.validate.domainListMatch(validator.result)
}

func (validator *validator) validateRegex() {
	validator.validate.regex(validator.result)
}

func (validator *validator) validateMx() {
	validator.validate.mx(validator.result)
}

func (validator *validator) validateSMTP() {
	validator.validate.smtp(validator.result)
}

func (validator *validator) run() *validatorResult {
	// preparing for running
	validator.usedValidations = []string{}

	// Whitelist/Blacklist validation
	validatorResult := validator.result
	validator.validateDomainListMatch()
	if !validatorResult.Success || !validator.isPassFromDomainListMatch {
		return validatorResult
	}
	// define validation flow
	switch validatorResult.ValidationType {
	case ValidationTypeRegex:
		validator.validateRegex()
	case ValidationTypeMx:
		validator.validateMx()
	case ValidationTypeSMTP:
		validator.validateSMTP()
	}
	return validatorResult
}

func (validator *validator) addUsedValidationType(validationType string) {
	validator.usedValidations = append(validator.usedValidations, validationType)
}
