package truemail

import "github.com/stretchr/testify/mock"

// Testing mocks

// validationDomainListMatch structure mock
type validationDomainListMatchMock struct {
	mock.Mock
}

func (validation *validationDomainListMatchMock) check(result *ValidatorResult) *ValidatorResult {
	args := validation.Called(result)
	return args.Get(0).(*ValidatorResult)
}

// validationRegex structure mock
type validationRegexMock struct {
	mock.Mock
}

func (validation *validationRegexMock) check(result *ValidatorResult) *ValidatorResult {
	args := validation.Called(result)
	return args.Get(0).(*ValidatorResult)
}

// validationMx structure mock
type validationMxMock struct {
	mock.Mock
}

func (validation *validationMxMock) check(result *ValidatorResult) *ValidatorResult {
	args := validation.Called(result)
	return args.Get(0).(*ValidatorResult)
}

// validationMxBlacklistMock structure mock
type validationMxBlacklistMock struct {
	mock.Mock
}

func (validation *validationMxBlacklistMock) check(result *ValidatorResult) *ValidatorResult {
	args := validation.Called(result)
	return args.Get(0).(*ValidatorResult)
}

// validationSmtpMock structure mock
type validationSmtpMock struct {
	mock.Mock
}

func (validation *validationSmtpMock) check(result *ValidatorResult) *ValidatorResult {
	args := validation.Called(result)
	return args.Get(0).(*ValidatorResult)
}

// dnsResolverMock structure mock
type dnsResolverMock struct {
	mock.Mock
}

func (resolver *dnsResolverMock) aRecord(hostName string) (string, error) {
	args := resolver.Called(hostName)
	return args.String(0), args.Error(1)
}

func (resolver *dnsResolverMock) aRecords(hostName string) ([]string, error) {
	args := resolver.Called(hostName)
	return args.Get(0).([]string), args.Error(1)
}

func (resolver *dnsResolverMock) cnameRecord(hostName string) (string, error) {
	args := resolver.Called(hostName)
	return args.String(0), args.Error(1)
}

func (resolver *dnsResolverMock) mxRecords(hostName string) ([]uint16, []string, error) {
	args := resolver.Called(hostName)
	return args.Get(0).([]uint16), args.Get(1).([]string), args.Error(2)
}

func (resolver *dnsResolverMock) ptrRecords(hostName string) ([]string, error) {
	args := resolver.Called(hostName)
	return args.Get(0).([]string), args.Error(1)
}

// smtpClientMock structure mock
type smtpClientMock struct {
	mock.Mock
}

func (client *smtpClientMock) sessionError() *SmtpClientError {
	return client.Called().Get(0).(*SmtpClientError)
}

func (client *smtpClientMock) runSession() bool {
	return client.Called().Bool(0)
}

// smtpBuilderMock structure mock
type smtpBuilderMock struct {
	mock.Mock
}

func (builder *smtpBuilderMock) newSmtpRequest(attempts int, targetEmail, targetHostAddress string, configuration *Configuration) *SmtpRequest {
	args := builder.Called(attempts, targetEmail, targetHostAddress, configuration)
	return args.Get(0).(*SmtpRequest)
}

func (builder *smtpBuilderMock) newSmtpClient(configuration *SmtpRequestConfiguration) client {
	args := builder.Called(configuration)
	return args.Get(0).(client)
}
