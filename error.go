package truemail

import "net"

type validationError struct {
	isDnsNotFound, isNullMxFound bool
	err                          error
}

func (customError *validationError) Error() string {
	return customError.err.Error()
}

func wrapNullMxError(err error) *validationError {
	return &validationError{isNullMxFound: true, err: err}
}

func wrapDnsError(err error) *validationError {
	e, ok := err.(*net.DNSError)
	if ok && e.IsNotFound {
		return &validationError{isDnsNotFound: true, err: err}
	}

	return &validationError{err: err}
}
