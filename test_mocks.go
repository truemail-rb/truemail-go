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
