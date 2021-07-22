package truemail

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigurationAttrAssignDefaultValues(t *testing.T) {
	t.Run("when created ConfigurationAttr structure with default field values", func(t *testing.T) {
		configurationAttr := new(ConfigurationAttr)
		configurationAttr.assignDefaultValues()

		assert.Equal(t, ValidationTypeDefault, configurationAttr.validationTypeDefault)
		assert.Equal(t, RegexEmailPattern, configurationAttr.emailPattern)
		assert.Equal(t, RegexSMTPErrorBodyPattern, configurationAttr.smtpErrorBodyPattern)
		assert.Equal(t, DefaultConnectionTimeout, configurationAttr.connectionTimeout)
		assert.Equal(t, DefaultResponseTimeout, configurationAttr.responseTimeout)
		assert.Equal(t, DefaultConnectionAttempts, configurationAttr.connectionAttempts)
	})

	t.Run("when created ConfigurationAttr structure with custom field values", func(t *testing.T) {
		validationTypeDefault, emailPattern, smtpErrorBodyPattern := "1", "2", "3"
		connectionTimeout, responseTimeout, connectionAttempts := 1, 2, 3
		configurationAttr := ConfigurationAttr{
			validationTypeDefault: validationTypeDefault,
			emailPattern:          emailPattern,
			smtpErrorBodyPattern:  smtpErrorBodyPattern,
			connectionTimeout:     connectionTimeout,
			responseTimeout:       responseTimeout,
			connectionAttempts:    connectionAttempts,
		}
		configurationAttr.assignDefaultValues()

		assert.Equal(t, validationTypeDefault, configurationAttr.validationTypeDefault)
		assert.Equal(t, emailPattern, configurationAttr.emailPattern)
		assert.Equal(t, smtpErrorBodyPattern, configurationAttr.smtpErrorBodyPattern)
		assert.Equal(t, connectionTimeout, configurationAttr.connectionTimeout)
		assert.Equal(t, responseTimeout, configurationAttr.responseTimeout)
		assert.Equal(t, connectionAttempts, configurationAttr.connectionAttempts)
	})
}

func TestConfigurationAttrValidate(t *testing.T) {
	t.Run("invalid verifier email", func(t *testing.T) {
		configurationAttr := ConfigurationAttr{verifierEmail: "email@domain"}
		errorMessage := fmt.Sprintf("%v is invalid verifier email", configurationAttr.verifierEmail)

		assert.EqualError(t, configurationAttr.validate(), errorMessage)
	})

	t.Run("invalid verifier domain", func(t *testing.T) {
		configurationAttr := ConfigurationAttr{
			verifierEmail:  randomEmail(),
			verifierDomain: "invalid_domain",
		}
		errorMessage := fmt.Sprintf("%v is invalid verifier domain", configurationAttr.verifierDomain)

		assert.EqualError(t, configurationAttr.validate(), errorMessage)
	})

	t.Run("invalid default validation type", func(t *testing.T) {
		configurationAttr := ConfigurationAttr{verifierEmail: randomEmail(), validationTypeDefault: "invalid validation type"}
		errorMessage := fmt.Sprintf("%v is invalid default validation type, use one of these: [regex mx mx_blacklist smtp]", configurationAttr.validationTypeDefault)

		assert.EqualError(t, configurationAttr.validate(), errorMessage)
	})

	t.Run("invalid connection timeout", func(t *testing.T) {
		configurationAttr := ConfigurationAttr{
			verifierEmail:         randomEmail(),
			validationTypeDefault: randomValidationType(),
			connectionTimeout:     randomNegativeNumber(),
		}
		errorMessage := fmt.Sprintf("%v should be a positive integer", configurationAttr.connectionTimeout)

		assert.EqualError(t, configurationAttr.validate(), errorMessage)
	})

	t.Run("invalid response timeout", func(t *testing.T) {
		configurationAttr := ConfigurationAttr{
			verifierEmail:         randomEmail(),
			validationTypeDefault: randomValidationType(),
			connectionTimeout:     randomPositiveNumber(),
			responseTimeout:       randomNegativeNumber(),
		}
		errorMessage := fmt.Sprintf("%v should be a positive integer", configurationAttr.responseTimeout)

		assert.EqualError(t, configurationAttr.validate(), errorMessage)
	})

	t.Run("invalid connection attempts", func(t *testing.T) {
		configurationAttr := ConfigurationAttr{
			verifierEmail:         randomEmail(),
			validationTypeDefault: randomValidationType(),
			connectionTimeout:     randomPositiveNumber(),
			responseTimeout:       randomPositiveNumber(),
			connectionAttempts:    randomNegativeNumber(),
		}
		errorMessage := fmt.Sprintf("%v should be a positive integer", configurationAttr.connectionAttempts)

		assert.EqualError(t, configurationAttr.validate(), errorMessage)
	})

	t.Run("invalid whitelisted domains", func(t *testing.T) {
		configurationAttr := ConfigurationAttr{
			verifierEmail:         randomEmail(),
			validationTypeDefault: randomValidationType(),
			connectionTimeout:     randomPositiveNumber(),
			responseTimeout:       randomPositiveNumber(),
			connectionAttempts:    randomPositiveNumber(),
			whitelistedDomains:    []string{randomDomain(), "a"},
		}
		errorMessage := fmt.Sprintf("%v is invalid domain name", configurationAttr.whitelistedDomains[1])

		assert.EqualError(t, configurationAttr.validate(), errorMessage)
	})

	t.Run("invalid blacklisted domains", func(t *testing.T) {
		configurationAttr := ConfigurationAttr{
			verifierEmail:         randomEmail(),
			validationTypeDefault: randomValidationType(),
			connectionTimeout:     randomPositiveNumber(),
			responseTimeout:       randomPositiveNumber(),
			connectionAttempts:    randomPositiveNumber(),
			blacklistedDomains:    []string{randomDomain(), "b"},
		}
		errorMessage := fmt.Sprintf("%v is invalid domain name", configurationAttr.blacklistedDomains[1])

		assert.EqualError(t, configurationAttr.validate(), errorMessage)
	})

	t.Run("invalid blacklisted mx ip address", func(t *testing.T) {
		configurationAttr := ConfigurationAttr{
			verifierEmail:            randomEmail(),
			validationTypeDefault:    randomValidationType(),
			connectionTimeout:        randomPositiveNumber(),
			responseTimeout:          randomPositiveNumber(),
			connectionAttempts:       randomPositiveNumber(),
			blacklistedMxIpAddresses: []string{randomIpAddress(), "1.1.1.256:65536"},
		}
		errorMessage := fmt.Sprintf("%v is invalid ip address", configurationAttr.blacklistedMxIpAddresses[1])

		assert.EqualError(t, configurationAttr.validate(), errorMessage)
	})

	t.Run("invalid dns, wrong ip address", func(t *testing.T) {
		configurationAttr := ConfigurationAttr{
			verifierEmail:         randomEmail(),
			validationTypeDefault: randomValidationType(),
			connectionTimeout:     randomPositiveNumber(),
			responseTimeout:       randomPositiveNumber(),
			connectionAttempts:    randomPositiveNumber(),
			dns:                   "1.1.1.256",
		}
		errorMessage := fmt.Sprintf("%v is invalid dns server", configurationAttr.dns)

		assert.EqualError(t, configurationAttr.validate(), errorMessage)
	})

	t.Run("invalid dns, wrong port number", func(t *testing.T) {
		configurationAttr := ConfigurationAttr{
			verifierEmail:         randomEmail(),
			validationTypeDefault: randomValidationType(),
			connectionTimeout:     randomPositiveNumber(),
			responseTimeout:       randomPositiveNumber(),
			connectionAttempts:    randomPositiveNumber(),
			dns:                   "1.1.1.255:65536",
		}
		errorMessage := fmt.Sprintf("%v is invalid dns server", configurationAttr.dns)

		assert.EqualError(t, configurationAttr.validate(), errorMessage)
	})

	t.Run("invalid dns, wrong ip address and port number", func(t *testing.T) {
		configurationAttr := ConfigurationAttr{
			verifierEmail:         randomEmail(),
			validationTypeDefault: randomValidationType(),
			connectionTimeout:     randomPositiveNumber(),
			responseTimeout:       randomPositiveNumber(),
			connectionAttempts:    randomPositiveNumber(),
			dns:                   "1.1.1.256:65536",
		}
		errorMessage := fmt.Sprintf("%v is invalid dns server", configurationAttr.dns)

		assert.EqualError(t, configurationAttr.validate(), errorMessage)
	})

	t.Run("invalid validation type by domain, wrong domain", func(t *testing.T) {
		invalidDomain := "inavlid domain"
		configurationAttr := ConfigurationAttr{
			verifierEmail:          randomEmail(),
			validationTypeDefault:  randomValidationType(),
			connectionTimeout:      randomPositiveNumber(),
			responseTimeout:        randomPositiveNumber(),
			connectionAttempts:     randomPositiveNumber(),
			validationTypeByDomain: map[string]string{randomDomain(): "regex", invalidDomain: "wrong_type"},
		}
		errorMessage := fmt.Sprintf("%v is invalid domain name", invalidDomain)

		assert.EqualError(t, configurationAttr.validate(), errorMessage)
	})

	t.Run("invalid validation type by domain, wrong validation type", func(t *testing.T) {
		invalidType := "inavlid validation type"
		configurationAttr := ConfigurationAttr{
			verifierEmail:          randomEmail(),
			validationTypeDefault:  randomValidationType(),
			connectionTimeout:      randomPositiveNumber(),
			responseTimeout:        randomPositiveNumber(),
			connectionAttempts:     randomPositiveNumber(),
			validationTypeByDomain: map[string]string{randomDomain(): "regex", randomDomain(): invalidType},
		}
		errorMessage := fmt.Sprintf("%v is invalid default validation type, use one of these: [regex mx mx_blacklist smtp]", invalidType)

		assert.EqualError(t, configurationAttr.validate(), errorMessage)
	})

	t.Run("invalid email pattern", func(t *testing.T) {
		configurationAttr := ConfigurationAttr{
			verifierEmail:         randomEmail(),
			validationTypeDefault: randomValidationType(),
			connectionTimeout:     randomPositiveNumber(),
			responseTimeout:       randomPositiveNumber(),
			connectionAttempts:    randomPositiveNumber(),
			emailPattern:          `\K`,
		}
		errorMessage := fmt.Sprintf("error parsing regexp: invalid escape sequence: `%v`", configurationAttr.emailPattern)

		assert.EqualError(t, configurationAttr.validate(), errorMessage)
	})

	t.Run("invalid smtp error body pattern", func(t *testing.T) {
		configurationAttr := ConfigurationAttr{
			verifierEmail:         randomEmail(),
			validationTypeDefault: randomValidationType(),
			connectionTimeout:     randomPositiveNumber(),
			responseTimeout:       randomPositiveNumber(),
			connectionAttempts:    randomPositiveNumber(),
			smtpErrorBodyPattern:  `\K`,
		}
		errorMessage := fmt.Sprintf("error parsing regexp: invalid escape sequence: `%v`", configurationAttr.smtpErrorBodyPattern)

		assert.EqualError(t, configurationAttr.validate(), errorMessage)
	})

	t.Run("coerces to special format/types", func(t *testing.T) {
		randomIpAddress, randomIpAddressWithDefaultPortNumber := randomDnsServerWithDefaultPortNumber()
		regexPatternFirst, regexPatternSecond := "1", "2"
		configurationAttr := ConfigurationAttr{
			verifierEmail:         randomEmail(),
			validationTypeDefault: randomValidationType(),
			connectionTimeout:     randomPositiveNumber(),
			responseTimeout:       randomPositiveNumber(),
			connectionAttempts:    randomPositiveNumber(),
			dns:                   randomIpAddress,
			emailPattern:          regexPatternFirst,
			smtpErrorBodyPattern:  regexPatternSecond,
		}
		regexObjectFirst, _ := newRegex(regexPatternFirst)
		regexObjectSecond, _ := newRegex(regexPatternSecond)

		assert.NoError(t, configurationAttr.validate())
		assert.Equal(t, randomIpAddressWithDefaultPortNumber, configurationAttr.dns)
		assert.EqualValues(t, regexObjectFirst, configurationAttr.regexEmail)
		assert.EqualValues(t, regexObjectSecond, configurationAttr.regexSMTPErrorBody)
	})
}

func TestConfigurationAttrValidateVerifierEmail(t *testing.T) {
	t.Run("valid verifier email", func(t *testing.T) {
		assert.NoError(t, new(ConfigurationAttr).validateVerifierEmail("el+niño!@mañana.es"))
	})

	t.Run("invalid verifier email", func(t *testing.T) {
		invalidEmail := "email"
		errorMessage := fmt.Sprintf("%s is invalid verifier email", invalidEmail)

		assert.EqualError(t, new(ConfigurationAttr).validateVerifierEmail(invalidEmail), errorMessage)
	})
}

func TestConfigurationAttrValidateVerifierDomain(t *testing.T) {
	t.Run("valid verifier domain", func(t *testing.T) {
		validDomain := "mañana.es"
		domain, err := new(ConfigurationAttr).validateVerifierDomain(validDomain)

		assert.Equal(t, validDomain, domain)
		assert.NoError(t, err)
	})

	t.Run("invalid verifier domain", func(t *testing.T) {
		invalidDomain := "domain"
		errorMessage := fmt.Sprintf("%s is invalid verifier domain", invalidDomain)
		domain, err := new(ConfigurationAttr).validateVerifierDomain(invalidDomain)

		assert.Equal(t, invalidDomain, domain)
		assert.EqualError(t, err, errorMessage)
	})
}

func TestConfigurationAttrBuildVerifierDomain(t *testing.T) {
	verifierEmail, expectedDomain := pairRandomEmailDomain()

	t.Run("valid verifier domain", func(t *testing.T) {
		validDomain := "mañana.es"
		domain, err := new(ConfigurationAttr).buildVerifierDomain(verifierEmail, validDomain)

		assert.Equal(t, validDomain, domain)
		assert.NoError(t, err)
	})

	t.Run("empty verifier domain", func(t *testing.T) {
		actualDomain, err := new(ConfigurationAttr).buildVerifierDomain(verifierEmail, EmptyString)

		assert.Equal(t, expectedDomain, actualDomain)
		assert.NoError(t, err)
	})

	t.Run("invalid verifier domain", func(t *testing.T) {
		invalidDomain := "domain"
		errorMessage := fmt.Sprintf("%s is invalid verifier domain", invalidDomain)
		domain, err := new(ConfigurationAttr).buildVerifierDomain(verifierEmail, invalidDomain)

		assert.Equal(t, invalidDomain, domain)
		assert.EqualError(t, err, errorMessage)
	})
}

func TestConfigurationAttrValidateValidationTypeDefaultContext(t *testing.T) {
	for _, validType := range availableValidationTypes() {
		t.Run("valid validation type", func(t *testing.T) {
			assert.NoError(t, new(ConfigurationAttr).validateValidationTypeDefaultContext(validType))
		})
	}

	t.Run("invalid validation type", func(t *testing.T) {
		invalidType := "invalid type"
		errorMessage := fmt.Sprintf("%s is invalid default validation type, use one of these: [regex mx mx_blacklist smtp]", invalidType)

		assert.EqualError(t, new(ConfigurationAttr).validateValidationTypeDefaultContext(invalidType), errorMessage)
	})
}

func TestConfigurationAttrValidateIntegerPositive(t *testing.T) {
	t.Run("valid positive integer", func(t *testing.T) {
		assert.NoError(t, new(ConfigurationAttr).validateIntegerPositive(42))
	})

	t.Run("invalid positive integer", func(t *testing.T) {
		notPositiveInteger := -42
		errorMessage := fmt.Sprintf("%v should be a positive integer", notPositiveInteger)

		assert.EqualError(t, new(ConfigurationAttr).validateIntegerPositive(notPositiveInteger), errorMessage)
	})
}

func TestConfigurationAttrValidateDomainContext(t *testing.T) {
	t.Run("valid domain", func(t *testing.T) {
		assert.NoError(t, new(ConfigurationAttr).validateDomainContext(randomDomain()))
	})

	t.Run("invalid domain", func(t *testing.T) {
		invalidDomain := "wrong.d"
		errorMessage := fmt.Sprintf("%s is invalid domain name", invalidDomain)

		assert.EqualError(t, new(ConfigurationAttr).validateDomainContext(invalidDomain), errorMessage)
	})
}

func TestConfigurationAttrValidateDomainsContext(t *testing.T) {
	t.Run("empty domains", func(t *testing.T) {
		assert.NoError(t, new(ConfigurationAttr).validateDomainsContext([]string{}))
	})

	t.Run("valid domains", func(t *testing.T) {
		assert.NoError(t, new(ConfigurationAttr).validateDomainsContext([]string{randomDomain(), "mañana.es"}))
	})

	t.Run("included invalid domain", func(t *testing.T) {
		invalidDomain := "wrong.d"
		domains := []string{randomDomain(), invalidDomain, "wrong.d2"}
		errorMessage := fmt.Sprintf("%s is invalid domain name", invalidDomain)

		assert.EqualError(t, new(ConfigurationAttr).validateDomainsContext(domains), errorMessage)
	})
}

func TestConfigurationAttrValidateIpAddressContext(t *testing.T) {
	t.Run("valid ip address", func(t *testing.T) {
		assert.NoError(t, new(ConfigurationAttr).validateIpAddressContext(randomIpAddress()))
	})

	invalidIpAddresses := []string{"10.300.0.256", "11.287.0.1", "172.1600.0.0", "-0.1.1.1", "8.08.8.8", "192.168.0.255a", "0.00.0.42"}

	for _, invalidIpAddress := range invalidIpAddresses {
		t.Run("invalid ip address", func(t *testing.T) {
			errorMessage := fmt.Sprintf("%s is invalid ip address", invalidIpAddress)

			assert.EqualError(t, new(ConfigurationAttr).validateIpAddressContext(invalidIpAddress), errorMessage)
		})
	}
}

func TestConfigurationAttrValidateIpAddressesContext(t *testing.T) {
	t.Run("empty ip addresses", func(t *testing.T) {
		assert.NoError(t, new(ConfigurationAttr).validateIpAddressesContext([]string{}))
	})

	t.Run("valid ip addresses", func(t *testing.T) {
		assert.NoError(t, new(ConfigurationAttr).validateIpAddressesContext([]string{randomIpAddress(), randomIpAddress()}))
	})

	t.Run("included invalid ip addresses", func(t *testing.T) {
		invalidIpAddress := "not_ip_address"
		ipAddresses := []string{randomIpAddress(), invalidIpAddress}
		errorMessage := fmt.Sprintf("%s is invalid ip address", invalidIpAddress)

		assert.EqualError(t, new(ConfigurationAttr).validateIpAddressesContext(ipAddresses), errorMessage)
	})
}

func TestConfigurationAttrValidateDNSServerContext(t *testing.T) {
	t.Run("valid dns server ip without port number", func(t *testing.T) {
		assert.NoError(t, new(ConfigurationAttr).validateDNSServerContext(randomIpAddress()))
	})

	t.Run("valid dns server ip with port number", func(t *testing.T) {
		assert.NoError(t, new(ConfigurationAttr).validateDNSServerContext(randomIpAddress()+":65507"))
	})

	t.Run("invalid dns server ip without port number", func(t *testing.T) {
		invalidDNSServer := "1.1.1.256"
		errorMessage := fmt.Sprintf("%s is invalid dns server", invalidDNSServer)

		assert.EqualError(t, new(ConfigurationAttr).validateDNSServerContext(invalidDNSServer), errorMessage)
	})

	t.Run("valid dns server ip with invalid port number", func(t *testing.T) {
		invalidDNSServer := "1.1.1.1:65536"
		errorMessage := fmt.Sprintf("%s is invalid dns server", invalidDNSServer)

		assert.EqualError(t, new(ConfigurationAttr).validateDNSServerContext(invalidDNSServer), errorMessage)
	})

	t.Run("invalid dns server ip with invalid port number", func(t *testing.T) {
		invalidDNSServer := "256.256.256.256:0"
		errorMessage := fmt.Sprintf("%s is invalid dns server", invalidDNSServer)

		assert.EqualError(t, new(ConfigurationAttr).validateDNSServerContext(invalidDNSServer), errorMessage)
	})
}

func TestConfigurationAttrValidateTypeByDomainContext(t *testing.T) {
	t.Run("empty dictionary", func(t *testing.T) {
		assert.NoError(t, new(ConfigurationAttr).validateTypeByDomainContext(map[string]string{}))
	})

	t.Run("valid dictionary", func(t *testing.T) {
		validTypesByDomains := map[string]string{randomDomain(): "regex", randomDomain(): "mx", randomDomain(): "smtp"}

		assert.NoError(t, new(ConfigurationAttr).validateTypeByDomainContext(validTypesByDomains))
	})

	for _, validationType := range []string{"regex", "invalid validation type"} {
		t.Run("included invalid domain", func(t *testing.T) {
			invalidDomain := "wrong.d"
			typesByDomains := map[string]string{invalidDomain: validationType}
			errorMessage := fmt.Sprintf("%s is invalid domain name", invalidDomain)

			assert.EqualError(t, new(ConfigurationAttr).validateTypeByDomainContext(typesByDomains), errorMessage)
		})
	}

	t.Run("included invalid validation type", func(t *testing.T) {
		wrongType := "wrong validation type"
		typesByDomains := map[string]string{randomDomain(): wrongType}
		errorMessage := fmt.Sprintf("%s is invalid default validation type, use one of these: [regex mx mx_blacklist smtp]", wrongType)

		assert.EqualError(t, new(ConfigurationAttr).validateTypeByDomainContext(typesByDomains), errorMessage)
	})
}

func TestConfigurationAttrFormatDns(t *testing.T) {
	t.Run("when DNS port not specified", func(t *testing.T) {
		randomIpAddress, randomIpAddressWithDefaultPortNumber := randomDnsServerWithDefaultPortNumber()

		assert.Equal(t, randomIpAddressWithDefaultPortNumber, new(ConfigurationAttr).formatDns(randomIpAddress))
	})

	t.Run("when DNS port specified", func(t *testing.T) {
		dnsGateway := randomIpAddress() + ":5300"

		assert.Equal(t, dnsGateway, new(ConfigurationAttr).formatDns(dnsGateway))
	})
}
