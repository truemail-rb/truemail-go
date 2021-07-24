package truemail

import (
	"fmt"
	"net"
	"strconv"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/foxcpp/go-mockdns"
)

// Testing helpers

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

func randomIp6Address() string {
	gofakeit.Seed(0)
	return gofakeit.IPv6Address()
}

func randomPositiveNumber() int {
	gofakeit.Seed(0)
	return gofakeit.Number(1, 42)
}

func randomNegativeNumber() int {
	gofakeit.Seed(0)
	return gofakeit.Number(-42, 0)
}

func randomPortNumber() int {
	gofakeit.Seed(0)
	return gofakeit.Number(1, 65535)
}

func randomDnsServer() string {
	gofakeit.Seed(0)
	return randomIpAddress() + ":" + strconv.Itoa(randomPortNumber())
}

func randomDnsServerWithDefaultPortNumber() (string, string) {
	ipAddress := randomIpAddress()
	return ipAddress, ipAddress + ":" + DefaultDnsPort
}

func createConfiguration() *configuration {
	configuration, _ := NewConfiguration(ConfigurationAttr{verifierEmail: randomEmail()})
	return configuration
}

func createValidatorResult(email string, configuration *configuration, options ...string) *validatorResult {
	validationType, _ := variadicValidationType(options, configuration.ValidationTypeDefault)
	validatorResult := &validatorResult{Email: email, Configuration: configuration, ValidationType: validationType}
	return validatorResult
}

func createSuccessfulValidatorResult(email string, configuration *configuration) *validatorResult {
	return &validatorResult{Email: email, Configuration: configuration, Success: true}
}

func randomValidationType() string {
	gofakeit.Seed(0)
	availableValidationTypes := []string{ValidationTypeRegex, ValidationTypeMx, ValidationTypeSMTP}
	index := gofakeit.Number(0, len(availableValidationTypes)-1)
	return availableValidationTypes[index]
}

func createValidator(email string, configuration *configuration, options ...string) *validator {
	validationType, _ := variadicValidationType(options, configuration.ValidationTypeDefault)
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

// Returns dnsResolver with mocked DNS records
func createDnsResolver(dnsRecords map[string]mockdns.Zone) *dnsResolver {
	return &dnsResolver{gateway: &mockdns.Resolver{Zones: dnsRecords}}
}

func createDnsResolverWithEpmtyRecords() *dnsResolver {
	return createDnsResolver(map[string]mockdns.Zone{})
}

func dnsErrorMessage(hostname string) string {
	return fmt.Sprintf("lookup %s on 127.0.0.1:53: no such host", hostname)
}

func runMockDnsServer(dnsRecords map[string]mockdns.Zone) string { // TODO: how to remove DNS request stdout dig log?
	srv, _ := mockdns.NewServer(dnsRecords, false)
	runningMockServerAddress := srv.LocalAddr().String()
	defer srv.Close()
	srv.PatchNet(net.DefaultResolver)
	defer mockdns.UnpatchNet(net.DefaultResolver)
	return runningMockServerAddress
}
