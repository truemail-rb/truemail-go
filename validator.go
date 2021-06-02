package truemail

import "fmt"

// ValidatorResult structure
type ValidatorResult struct {
	Success, SMTPDebug bool
	Email, Domain      string
	MailServers        []string
	Errors             map[string]string
	Configuration      *Configuration
}

func Validate(email string, configuration *Configuration, validationType string) (*ValidatorResult, error) {
	err := validateValidationTypeContext(validationType)
	if err != nil {
		return nil, err
	}

	validatorResult := newValidatorResult(email, configuration)

	// define validationType flow

	switch validationType {
	case ValidationTypeRegex:
		validateRegex(validatorResult)
	case ValidationTypeMx:
		validateMx(validateRegex(validatorResult))
	case ValidationTypeSMTP:
		validateSMTP(validateMx(validateRegex(validatorResult)))
	}

	return validatorResult, err
}

func validateValidationTypeContext(validationType string) error {
	if included(availableValidationTypes(), validationType) {
		return nil
	}
	return fmt.Errorf(
		"%s is invalid validation type, use one of these: %s",
		validationType,
		availableValidationTypes(),
	)
}

func newValidatorResult(email string, configuration *Configuration) *ValidatorResult {
	return &ValidatorResult{Email: email, Configuration: configuration}
}
