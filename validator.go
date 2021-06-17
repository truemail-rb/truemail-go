package truemail

import "fmt"

func Validate(email string, configuration *configuration, options ...string) (*validatorResult, error) {
	validationType, err := variadicValidationType(options)

	if err != nil {
		return nil, err
	}

	return newValidator(email, validationType, configuration).run(), err
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

// validator, structure with behaviour
type validator struct {
	result                    *validatorResult
	usedValidations           []string
	isPassFromDomainListMatch bool
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

func (validator *validator) validateRegex() {
	validateRegex(validator.result)
}

func (validator *validator) validateMx() {
	validateMx(validator.result)
}

func (validator *validator) validateSMTP() {
	validateSMTP(validator.result)
}

func (validator *validator) run() *validatorResult {
	// preparing for running
	validator.usedValidations = []string{}

	// Whitelist/Blacklist validation
	result := validator.result
	validateDomainListMatch(result)
	if !result.Success || !validator.isPassFromDomainListMatch {
		return result
	}
	// define validation flow
	switch result.ValidationType {
	case ValidationTypeRegex:
		validator.validateRegex()
	case ValidationTypeMx:
		validator.validateMx()
	case ValidationTypeSMTP:
		validator.validateSMTP()
	}
	return result
}

func (validator *validator) addUsedValidationType(validationType string) {
	validator.usedValidations = append(validator.usedValidations, validationType)
}
