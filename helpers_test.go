package truemail

import (
	"github.com/brianvoe/gofakeit/v6"
)

func createRandomEmail() string {
	gofakeit.Seed(0)
	return gofakeit.Email()
}

func createRandomDomain() string {
	gofakeit.Seed(0)
	return gofakeit.DomainName()
}

func createPairRandomEmailDomain() (string, string) {
	gofakeit.Seed(0)
	domain := createRandomDomain()
	email := gofakeit.Username() + "@" + domain
	return email, domain
}

func createConfiguration() *configuration {
	configuration, _ := NewConfiguration(ConfigurationAttr{verifierEmail: createRandomEmail()})
	return configuration
}

func createValidatorResult(email string, configuration *configuration, validationType ...string) *validatorResult {
	if len(validationType) == 0 {
		validationType = append(validationType, ValidationTypeDefault)
	}
	validatorResult := &validatorResult{Email: email, Configuration: configuration, ValidationType: validationType[0]}
	validatorResult.validator = &validator{result: validatorResult}
	return validatorResult
}

func createRandomValidationType() string {
	gofakeit.Seed(0)
	availableValidationTypes := []string{ValidationTypeRegex, ValidationTypeMx, ValidationTypeSMTP}
	index := gofakeit.Number(0, len(availableValidationTypes)-1)
	return availableValidationTypes[index]
}

func usedValidationsByType(validationType string) []string {
	return map[string][]string{
		ValidationTypeRegex: {ValidationTypeRegex},
		ValidationTypeMx:    {ValidationTypeRegex, ValidationTypeMx},
		ValidationTypeSMTP:  {ValidationTypeRegex, ValidationTypeMx, ValidationTypeSMTP},
	}[validationType]
}
