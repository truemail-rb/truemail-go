package truemail

import "fmt"

func Validate(validationAttr ValidationAttr) (*validatorResult, error) {
	validationType, validationConfiguration := validationAttr.validationType, validationAttr.configuration

	if validationType == "" {
		validationType = ValidationTypeDefault
	}

	err := validateValidationTypeContext(validationType)
	if err != nil {
		return nil, err
	}

	return newValidator(validationAttr.email, validationType, validationConfiguration).run(), err
}

// ValidationAttr kwargs for validator enrty point
type ValidationAttr struct {
	email, validationType string
	configuration         *configuration
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
