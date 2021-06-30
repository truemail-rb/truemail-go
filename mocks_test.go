package truemail

import "github.com/stretchr/testify/mock"

// validation structure mock
type validationMock struct {
	mock.Mock
}

// validation structure mock methods

func (validation *validationMock) domainListMatch(result *validatorResult) *validatorResult {
	args := validation.Called(result)
	return args.Get(0).(*validatorResult)
}

func (validation *validationMock) regex(result *validatorResult) *validatorResult {
	args := validation.Called(result)
	return args.Get(0).(*validatorResult)
}

func (validation *validationMock) mx(result *validatorResult) *validatorResult {
	args := validation.Called(result)
	return args.Get(0).(*validatorResult)
}

func (validation *validationMock) smtp(result *validatorResult) *validatorResult {
	args := validation.Called(result)
	return args.Get(0).(*validatorResult)
}
