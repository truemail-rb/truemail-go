package truemail

import (
	"github.com/brianvoe/gofakeit/v6"
)

func randomEmail() string {
	gofakeit.Seed(0)
	return gofakeit.Email()
}

func randomDomain() string {
	gofakeit.Seed(0)
	return gofakeit.DomainName()
}

func pairRandomEmailDomain() (string, string) {
	gofakeit.Seed(0)
	domain := randomDomain()
	email := gofakeit.Username() + "@" + domain
	return email, domain
}

func randomIpAddress() string {
	gofakeit.Seed(0)
	return gofakeit.IPv4Address()
}

func createConfiguration() *configuration {
	configuration, _ := NewConfiguration(ConfigurationAttr{verifierEmail: randomEmail()})
	return configuration
}

func createValidatorResult(email string, configuration *configuration, options ...string) *validatorResult {
	validationType, _ := variadicValidationType(options)
	validatorResult := &validatorResult{Email: email, Configuration: configuration, ValidationType: validationType}
	validatorResult.validator = &validator{result: validatorResult}
	return validatorResult
}

func randomValidationType() string {
	gofakeit.Seed(0)
	availableValidationTypes := []string{ValidationTypeRegex, ValidationTypeMx, ValidationTypeSMTP}
	index := gofakeit.Number(0, len(availableValidationTypes)-1)
	return availableValidationTypes[index]
}

func createValidator(email string, configuration *configuration, options ...string) *validator {
	validationType, _ := variadicValidationType(options)
	return newValidator(email, validationType, configuration)
}

func createValidatorPassedFromDomainListMatch(email string, configuration *configuration, options ...string) *validator {
	validator := createValidator(email, configuration, options...)
	validator.isPassFromDomainListMatch = true
	validator.result.Success = true
	return validator
}

func usedValidationsByType(validationType string) []string {
	return map[string][]string{
		ValidationTypeRegex: {ValidationTypeRegex},
		ValidationTypeMx:    {ValidationTypeRegex, ValidationTypeMx},
		ValidationTypeSMTP:  {ValidationTypeRegex, ValidationTypeMx, ValidationTypeSMTP},
	}[validationType]
}

func runDomainListMatchValidation(email string, configuration *configuration, options ...string) *validatorResult {
	validator := createValidator(email, configuration, options...)
	validatorResult := validator.result
	return validator.validate.domainListMatch(validatorResult)
}
