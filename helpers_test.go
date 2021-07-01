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

func usedValidationsByType(validationType string) []string {
	return map[string][]string{
		ValidationTypeRegex:       {ValidationTypeRegex},
		ValidationTypeMx:          {ValidationTypeRegex, ValidationTypeMx},
		ValidationTypeMxBlacklist: {ValidationTypeRegex, ValidationTypeMx, ValidationTypeMxBlacklist},
		ValidationTypeSMTP:        {ValidationTypeRegex, ValidationTypeMx, ValidationTypeMxBlacklist, ValidationTypeSMTP},
	}[validationType]
}

func runDomainListMatchValidation(email string, configuration *configuration, options ...string) *validatorResult {
	validator := createValidator(email, configuration, options...)
	validatorResult := validator.result
	return validator.domainListMatch.check(validatorResult)
}

func doPassedFromDomainListMatch(validatorResult *validatorResult) {
	validatorResult.Success, validatorResult.isPassFromDomainListMatch = true, true
}

func failedValidatorResult() *validatorResult {
	return new(validatorResult)
}
