package truemail

import "github.com/stretchr/testify/mock"

// Testing mocks

// validationDomainListMatch structure mock
type validationDomainListMatchMock struct {
	mock.Mock
}

func (validation *validationDomainListMatchMock) check(result *validatorResult) *validatorResult {
	args := validation.Called(result)
	return args.Get(0).(*validatorResult)
}

// validationRegex structure mock
type validationRegexMock struct {
	mock.Mock
}

func (validation *validationRegexMock) check(result *validatorResult) *validatorResult {
	args := validation.Called(result)
	return args.Get(0).(*validatorResult)
}

// validationMx structure mock
type validationMxMock struct {
	mock.Mock
}

func (validation *validationMxMock) check(result *validatorResult) *validatorResult {
	args := validation.Called(result)
	return args.Get(0).(*validatorResult)
}

// validationMxBlacklistMock structure mock
type validationMxBlacklistMock struct {
	mock.Mock
}

func (validation *validationMxBlacklistMock) check(result *validatorResult) *validatorResult {
	args := validation.Called(result)
	return args.Get(0).(*validatorResult)
}

// validationSmtpMock structure mock
type validationSmtpMock struct {
	mock.Mock
}

func (validation *validationSmtpMock) check(result *validatorResult) *validatorResult {
	args := validation.Called(result)
	return args.Get(0).(*validatorResult)
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
