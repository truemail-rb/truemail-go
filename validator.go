package truemail

import "fmt"

// validatorResult structure
type validatorResult struct {
	Success, SMTPDebug            bool
	Email, Domain, ValidationType string
	MailServers                   []string
	Errors                        map[string]string
	Configuration                 *configuration
	isPassFromDomainListMatch     bool
}

// ValidationAttr kwargs for validator enrty point
type ValidationAttr struct {
	email, validationType string
	configuration         *configuration
}

func Validate(validationAttr ValidationAttr) (*validatorResult, error) {
	validationType, validationConfiguration := &validationAttr.validationType, validationAttr.configuration

	if *validationType == "" {
		*validationType = ValidationTypeDefault
	}

	err := validateValidationTypeContext(*validationType)
	if err != nil {
		return nil, err
	}

	validatorResult := newValidatorResult(validationAttr.email, validationConfiguration, *validationType)

	// define validationType flow

	// Whitelist/Blacklist validation
	validateDomainListMatch(validatorResult)
	if !validatorResult.Success || !validatorResult.isPassFromDomainListMatch {
		return validatorResult, err
	}

	// switch validationType {
	// case ValidationTypeRegex:
	// 	validateRegex(validatorResult)
	// case ValidationTypeMx:
	// 	validateMx(validateRegex(validatorResult))
	// case ValidationTypeSMTP:
	// 	validateSMTP(validateMx(validateRegex(validatorResult)))
	// }

	return validatorResult, err
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

func newValidatorResult(email string, configuration *configuration, validationType string) *validatorResult {
	return &validatorResult{Email: email, Configuration: configuration, ValidationType: validationType}
}

func addError(validatorResult *validatorResult, key, value string) *validatorResult {
	if validatorResult.Errors == nil {
		validatorResult.Errors = map[string]string{}
	}
	validatorResult.Errors[key] = value
	return validatorResult
}
