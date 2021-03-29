package truemail

import "github.com/brianvoe/gofakeit/v6"

func createRandomEmail() string {
	gofakeit.Seed(0)
	return gofakeit.Email()
}

func createConfiguration() *Configuration {
	configuration, _ := NewConfiguration(ConfigurationAttr{verifierEmail: createRandomEmail()})
	return configuration
}
