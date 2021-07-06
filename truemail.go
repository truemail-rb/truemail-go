package truemail

import "fmt"

func Validate(email string, configuration *configuration, options ...string) (*validatorResult, error) {
	validationType, err := variadicValidationType(options)

	if err != nil {
		return nil, err
	}

	return newValidator(email, validationType, configuration).run(), err
}

func IsValid(email string, configuration *configuration, options ...string) bool {
	validationType, err := variadicValidationType(options)

	if err != nil {
		return false
	}

	return newValidator(email, validationType, configuration).run().Success
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
