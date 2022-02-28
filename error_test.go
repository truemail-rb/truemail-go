package truemail

import (
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidationErrorError(t *testing.T) {
	t.Run("returns wrapped error message", func(t *testing.T) {
		errorMessage := "error message"
		customError := &validationError{err: fmt.Errorf(errorMessage)}

		assert.Equal(t, errorMessage, customError.Error())
	})
}

func TestWrapNullMxError(t *testing.T) {
	t.Run("wrappes error to validationError with isNullMxFound marker", func(t *testing.T) {
		errorMessage := "error message"
		err := wrapNullMxError(fmt.Errorf(errorMessage))

		assert.True(t, err.isNullMxFound)
		assert.Equal(t, errorMessage, err.Error())
	})
}

func TestWrapDnsError(t *testing.T) {
	hostname, server, errMessage := randomDomain(), localhostIPv4Address+":53", "no such host"

	t.Run("when DNSError not found error", func(t *testing.T) {
		err := wrapDnsError(&net.DNSError{Name: hostname, Server: server, Err: errMessage, IsNotFound: true})

		assert.True(t, err.isDnsNotFound)
		assert.Equal(t, dnsErrorMessage(hostname), err.Error())
	})

	t.Run("when other DNSError error", func(t *testing.T) {
		err := wrapDnsError(&net.DNSError{Name: hostname, Server: server, Err: errMessage, IsTimeout: true})

		assert.False(t, err.isDnsNotFound)
		assert.Equal(t, dnsErrorMessage(hostname), err.Error())
	})
}

func TestSmtpClientErrorError(t *testing.T) {
	t.Run("returns wrapped error message", func(t *testing.T) {
		errorMessage := "error message"
		customError := &smtpClientError{err: fmt.Errorf(errorMessage)}

		assert.Equal(t, errorMessage, customError.Error())
	})
}
