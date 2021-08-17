package truemail

import "net"

// Error wrapper
type validationError struct {
	isDnsNotFound, isNullMxFound bool
	err                          error
}

// error interface implementation
func (customError *validationError) Error() string {
	return customError.err.Error()
}

// Wrappes error in validationError with isNullMxFound: true
func wrapNullMxError(err error) *validationError {
	return &validationError{isNullMxFound: true, err: err}
}

// Wrappes DNSError in validationError with isDnsNotFound,
// that depends on DNSError context
func wrapDnsError(err error) *validationError {
	e, ok := err.(*net.DNSError)
	if ok && e.IsNotFound {
		return &validationError{isDnsNotFound: true, err: err}
	}

	return &validationError{err: err}
}
