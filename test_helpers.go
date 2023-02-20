package truemail

import (
	"fmt"
	"net"
	"strconv"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/foxcpp/go-mockdns"
	smtpmock "github.com/mocktools/go-smtp-mock/v2"
	"golang.org/x/net/idna"
)

// truemail test helpers

var localhostIPv4Address = "127.0.0.1"

func randomEmail() string {
	gofakeit.Seed(0)
	return gofakeit.Email()
}

func randomDomain() string {
	gofakeit.Seed(0)
	return gofakeit.DomainName()
}

func randomDnsHostName() string {
	gofakeit.Seed(0)
	return gofakeit.DomainName() + "."
}

func punycodeDomain(domain string) string {
	punycodeDomain, _ := idna.New().ToASCII(domain)
	return punycodeDomain
}

func toDnsHostName(domain string) string {
	return domain + "."
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
	return ipAddress, serverWithPortNumber(ipAddress, defaultDnsPort)
}

func createConfiguration() *Configuration {
	configuration, _ := NewConfiguration(ConfigurationAttr{VerifierEmail: randomEmail()})
	return configuration
}

func createValidatorResult(email string, configuration *Configuration, options ...string) *ValidatorResult {
	validationType, _ := variadicValidationType(options, configuration.ValidationTypeDefault)
	return &ValidatorResult{Email: email, Configuration: configuration, ValidationType: validationType}
}

func createSuccessfulValidatorResult(email string, configuration *Configuration) *ValidatorResult {
	return &ValidatorResult{Email: email, Domain: emailDomain(email), Configuration: copyConfigurationByPointer(configuration), Success: true}
}

func randomValidationType() string {
	gofakeit.Seed(0)
	availableValidationTypes := []string{validationTypeRegex, validationTypeMx, validationTypeSmtp}
	index := gofakeit.Number(0, len(availableValidationTypes)-1)
	return availableValidationTypes[index]
}

func createValidator(email string, configuration *Configuration, options ...string) *validator {
	validationType, _ := variadicValidationType(options, configuration.ValidationTypeDefault)
	return newValidator(email, validationType, configuration)
}

func usedValidationsByType(validationType string) []string {
	return map[string][]string{
		validationTypeRegex:       {validationTypeRegex},
		validationTypeMx:          {validationTypeRegex, validationTypeMx},
		validationTypeMxBlacklist: {validationTypeRegex, validationTypeMx, validationTypeMxBlacklist},
		validationTypeSmtp:        {validationTypeRegex, validationTypeMx, validationTypeMxBlacklist, validationTypeSmtp},
	}[validationType]
}

func runDomainListMatchValidation(email string, configuration *Configuration, options ...string) *ValidatorResult {
	validator := createValidator(email, configuration, options...)
	validatorResult := validator.result
	return validator.domainListMatchLayer.check(validatorResult)
}

func doPassedFromDomainListMatch(validatorResult *ValidatorResult) {
	validatorResult.Success, validatorResult.isPassFromDomainListMatch = true, true
}

func failedValidatorResult() *ValidatorResult {
	return new(ValidatorResult)
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

// Runs DNS mock server. Returns running mock server address
func runMockDnsServer(dnsRecords map[string]mockdns.Zone) string { // TODO: how to remove DNS request stdout dig log?
	srv, _ := mockdns.NewServer(dnsRecords, false)
	runningMockServerAddress := srv.LocalAddr().String()
	defer srv.Close()
	srv.PatchNet(net.DefaultResolver)
	defer mockdns.UnpatchNet(net.DefaultResolver)
	return runningMockServerAddress
}

func createDnsNotFoundError() *validationError {
	return &validationError{isDnsNotFound: true, err: &net.DNSError{IsNotFound: true}}
}

func isDnsNotFoundError(err error) bool {
	e, ok := err.(*validationError)
	return ok && e.isDnsNotFound
}

func isNullMxError(err error) bool {
	e, ok := err.(*validationError)
	return ok && e.isNullMxFound
}

func startSmtpMock(config smtpmock.ConfigurationAttr) *smtpmock.Server {
	server := smtpmock.New(config)
	_ = server.Start()

	return server
}
