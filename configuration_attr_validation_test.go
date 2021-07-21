package truemail

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateVerifierEmail(t *testing.T) {
	t.Run("valid verifier email", func(t *testing.T) {
		assert.NoError(t, validateVerifierEmail("el+niño!@mañana.es"))
	})

	t.Run("invalid verifier email", func(t *testing.T) {
		invalidEmail := "email"
		errorMessage := fmt.Sprintf("%s is invalid verifier email", invalidEmail)

		assert.EqualError(t, validateVerifierEmail(invalidEmail), errorMessage)
	})
}

func TestValidateVerifierDomain(t *testing.T) {
	t.Run("valid verifier domain", func(t *testing.T) {
		validDomain := "mañana.es"
		domain, err := validateVerifierDomain(validDomain)

		assert.Equal(t, validDomain, domain)
		assert.NoError(t, err)
	})

	t.Run("invalid verifier domain", func(t *testing.T) {
		invalidDomain := "domain"
		errorMessage := fmt.Sprintf("%s is invalid verifier domain", invalidDomain)
		domain, err := validateVerifierDomain(invalidDomain)

		assert.Equal(t, invalidDomain, domain)
		assert.EqualError(t, err, errorMessage)
	})
}

func TestBuildVerifierDomain(t *testing.T) {
	verifierEmail, expectedDomain := pairRandomEmailDomain()

	t.Run("valid verifier domain", func(t *testing.T) {
		validDomain := "mañana.es"
		domain, err := buildVerifierDomain(verifierEmail, validDomain)

		assert.Equal(t, validDomain, domain)
		assert.NoError(t, err)
	})

	t.Run("empty verifier domain", func(t *testing.T) {
		actualDomain, err := buildVerifierDomain(verifierEmail, EmptyString)

		assert.Equal(t, expectedDomain, actualDomain)
		assert.NoError(t, err)
	})

	t.Run("invalid verifier domain", func(t *testing.T) {
		invalidDomain := "domain"
		errorMessage := fmt.Sprintf("%s is invalid verifier domain", invalidDomain)
		domain, err := buildVerifierDomain(verifierEmail, invalidDomain)

		assert.Equal(t, invalidDomain, domain)
		assert.EqualError(t, err, errorMessage)
	})
}

func TestAvailableValidationTypes(t *testing.T) {
	t.Run("slice of available validation types", func(t *testing.T) {
		assert.Equal(t, []string{"regex", "mx", "mx_blacklist", "smtp"}, availableValidationTypes())
	})
}

func TestValidateValidationTypeDefaultContext(t *testing.T) {
	for _, validType := range availableValidationTypes() {
		t.Run("valid validation type", func(t *testing.T) {
			assert.NoError(t, validateValidationTypeDefaultContext(validType))
		})
	}

	t.Run("invalid validation type", func(t *testing.T) {
		invalidType := "invalid type"
		errorMessage := fmt.Sprintf("%s is invalid default validation type, use one of these: [regex mx mx_blacklist smtp]", invalidType)

		assert.EqualError(t, validateValidationTypeDefaultContext(invalidType), errorMessage)
	})
}

func TestValidateIntegerPositive(t *testing.T) {
	t.Run("valid positive integer", func(t *testing.T) {
		assert.NoError(t, validateIntegerPositive(42))
	})

	t.Run("invalid positive integer", func(t *testing.T) {
		notPositiveInteger := -42
		errorMessage := fmt.Sprintf("%v should be a positive integer", notPositiveInteger)

		assert.EqualError(t, validateIntegerPositive(notPositiveInteger), errorMessage)
	})
}

func TestValidateDomainContext(t *testing.T) {
	t.Run("valid domain", func(t *testing.T) {
		assert.NoError(t, validateDomainContext(randomDomain()))
	})

	t.Run("invalid domain", func(t *testing.T) {
		invalidDomain := "wrong.d"
		errorMessage := fmt.Sprintf("%s is invalid domain name", invalidDomain)

		assert.EqualError(t, validateDomainContext(invalidDomain), errorMessage)
	})
}

func TestValidateDomainsContext(t *testing.T) {
	t.Run("empty domains", func(t *testing.T) {
		assert.NoError(t, validateDomainsContext([]string{}))
	})

	t.Run("valid domains", func(t *testing.T) {
		assert.NoError(t, validateDomainsContext([]string{randomDomain(), "mañana.es"}))
	})

	t.Run("included invalid domain", func(t *testing.T) {
		invalidDomain := "wrong.d"
		domains := []string{randomDomain(), invalidDomain, "wrong.d2"}
		errorMessage := fmt.Sprintf("%s is invalid domain name", invalidDomain)

		assert.EqualError(t, validateDomainsContext(domains), errorMessage)
	})
}

func TestValidateIpAddressContext(t *testing.T) {
	t.Run("valid ip address", func(t *testing.T) {
		assert.NoError(t, validateIpAddressContext(randomIpAddress()))
	})

	invalidIpAddresses := []string{"10.300.0.256", "11.287.0.1", "172.1600.0.0", "-0.1.1.1", "8.08.8.8", "192.168.0.255a", "0.00.0.42"}

	for _, invalidIpAddress := range invalidIpAddresses {
		t.Run("invalid ip address", func(t *testing.T) {
			errorMessage := fmt.Sprintf("%s is invalid ip address", invalidIpAddress)

			assert.EqualError(t, validateIpAddressContext(invalidIpAddress), errorMessage)
		})
	}
}

func TestValidateIpAddressesContext(t *testing.T) {
	t.Run("empty ip addresses", func(t *testing.T) {
		assert.NoError(t, validateIpAddressesContext([]string{}))
	})

	t.Run("valid ip addresses", func(t *testing.T) {
		assert.NoError(t, validateIpAddressesContext([]string{randomIpAddress(), randomIpAddress()}))
	})

	t.Run("included invalid ip addresses", func(t *testing.T) {
		invalidIpAddress := "not_ip_address"
		ipAddresses := []string{randomIpAddress(), invalidIpAddress}
		errorMessage := fmt.Sprintf("%s is invalid ip address", invalidIpAddress)

		assert.EqualError(t, validateIpAddressesContext(ipAddresses), errorMessage)
	})
}

func TestValidateDNSServerContext(t *testing.T) {
	t.Run("valid dns server ip without port number", func(t *testing.T) {
		assert.NoError(t, validateDNSServerContext(randomIpAddress()))
	})

	t.Run("valid dns server ip with port number", func(t *testing.T) {
		assert.NoError(t, validateDNSServerContext(randomIpAddress()+":65507"))
	})

	t.Run("invalid dns server ip without port number", func(t *testing.T) {
		invalidDNSServer := "1.1.1.256"
		errorMessage := fmt.Sprintf("%s is invalid dns server", invalidDNSServer)

		assert.EqualError(t, validateDNSServerContext(invalidDNSServer), errorMessage)
	})

	t.Run("valid dns server ip with invalid port number", func(t *testing.T) {
		invalidDNSServer := "1.1.1.1:65536"
		errorMessage := fmt.Sprintf("%s is invalid dns server", invalidDNSServer)

		assert.EqualError(t, validateDNSServerContext(invalidDNSServer), errorMessage)
	})

	t.Run("invalid dns server ip with invalid port number", func(t *testing.T) {
		invalidDNSServer := "256.256.256.256:0"
		errorMessage := fmt.Sprintf("%s is invalid dns server", invalidDNSServer)

		assert.EqualError(t, validateDNSServerContext(invalidDNSServer), errorMessage)
	})
}

func TestValidateTypeByDomainContext(t *testing.T) {
	t.Run("empty dictionary", func(t *testing.T) {
		assert.NoError(t, validateTypeByDomainContext(map[string]string{}))
	})

	t.Run("valid dictionary", func(t *testing.T) {
		validTypesByDomains := map[string]string{randomDomain(): "regex", randomDomain(): "mx", randomDomain(): "smtp"}

		assert.NoError(t, validateTypeByDomainContext(validTypesByDomains))
	})

	for _, validationType := range []string{"regex", "invalid validation type"} {
		t.Run("included invalid domain", func(t *testing.T) {
			invalidDomain := "wrong.d"
			typesByDomains := map[string]string{invalidDomain: validationType}
			errorMessage := fmt.Sprintf("%s is invalid domain name", invalidDomain)

			assert.EqualError(t, validateTypeByDomainContext(typesByDomains), errorMessage)
		})
	}

	t.Run("included invalid validation type", func(t *testing.T) {
		wrongType := "wrong validation type"
		typesByDomains := map[string]string{randomDomain(): wrongType}
		errorMessage := fmt.Sprintf("%s is invalid default validation type, use one of these: [regex mx mx_blacklist smtp]", wrongType)

		assert.EqualError(t, validateTypeByDomainContext(typesByDomains), errorMessage)
	})
}

func TestIsIncluded(t *testing.T) {
	var item string

	t.Run("item found in slice", func(t *testing.T) {
		assert.True(t, isIncluded([]string{item}, item))
	})

	t.Run("item not found in slice", func(t *testing.T) {
		assert.False(t, isIncluded([]string{}, item))
	})
}

func TestIsIntersected(t *testing.T) {
	var item string

	t.Run("item from target slice found in base slice", func(t *testing.T) {
		assert.True(t, isIntersected([]string{item}, []string{item}))
	})

	t.Run("item from target slice not found in base slice", func(t *testing.T) {
		assert.False(t, isIntersected([]string{}, []string{item}))
	})
}

func TestNewRegex(t *testing.T) {
	t.Run("valid regex pattern", func(t *testing.T) {
		regexPattern := EmptyString
		actualRegex, err := newRegex(regexPattern)
		expectedRegex, _ := regexp.Compile(regexPattern)

		assert.Equal(t, expectedRegex, actualRegex)
		assert.NoError(t, err)
	})

	t.Run("invalid regex pattern", func(t *testing.T) {
		actualRegex, err := newRegex(`\K`)

		assert.Nil(t, actualRegex)
		assert.Error(t, err)
	})
}

func TestMatchRegex(t *testing.T) {
	t.Run("valid regex pattern, matched string", func(t *testing.T) {
		assert.True(t, matchRegex(EmptyString, EmptyString))
	})

	t.Run("valid regex pattern, not matched string", func(t *testing.T) {
		assert.False(t, matchRegex("42", `\D+`))
	})

	t.Run("invalid regex pattern", func(t *testing.T) {
		assert.False(t, matchRegex(EmptyString, `\K`))
	})
}

func TestFormatDns(t *testing.T) {
	t.Run("when DNS port not specified", func(t *testing.T) {
		dnsGateway := randomIpAddress()

		assert.Equal(t, dnsGateway+":"+DefaultDnsPort, formatDns(dnsGateway))
	})

	t.Run("when DNS port specified", func(t *testing.T) {
		dnsGateway := randomIpAddress() + ":5300"

		assert.Equal(t, dnsGateway, formatDns(dnsGateway))
	})
}
