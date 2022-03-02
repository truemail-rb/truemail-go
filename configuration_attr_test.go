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

		assert.Equal(t, validationTypeDefault, configurationAttr.ValidationTypeDefault)
		assert.Equal(t, regexEmailPattern, configurationAttr.EmailPattern)
		assert.Equal(t, regexSMTPErrorBodyPattern, configurationAttr.SmtpErrorBodyPattern)
		assert.Equal(t, defaultConnectionTimeout, configurationAttr.ConnectionTimeout)
		assert.Equal(t, defaultResponseTimeout, configurationAttr.ResponseTimeout)
		assert.Equal(t, defaultConnectionAttempts, configurationAttr.ConnectionAttempts)
		assert.Equal(t, defaultSmtpPort, configurationAttr.SmtpPort)
	})

	t.Run("when created ConfigurationAttr structure with custom field values", func(t *testing.T) {
		ValidationTypeDefault, emailPattern, smtpErrorBodyPattern := "1", "2", "3"
		connectionTimeout, responseTimeout, connectionAttempts, smtpPort := 1, 2, 3, 4
		configurationAttr := ConfigurationAttr{
			ValidationTypeDefault: ValidationTypeDefault,
			EmailPattern:          emailPattern,
			SmtpErrorBodyPattern:  smtpErrorBodyPattern,
			ConnectionTimeout:     connectionTimeout,
			ResponseTimeout:       responseTimeout,
			ConnectionAttempts:    connectionAttempts,
			SmtpPort:              smtpPort,
		}
		configurationAttr.assignDefaultValues()

		assert.Equal(t, ValidationTypeDefault, configurationAttr.ValidationTypeDefault)
		assert.Equal(t, emailPattern, configurationAttr.EmailPattern)
		assert.Equal(t, smtpErrorBodyPattern, configurationAttr.SmtpErrorBodyPattern)
		assert.Equal(t, connectionTimeout, configurationAttr.ConnectionTimeout)
		assert.Equal(t, responseTimeout, configurationAttr.ResponseTimeout)
		assert.Equal(t, connectionAttempts, configurationAttr.ConnectionAttempts)
		assert.Equal(t, smtpPort, configurationAttr.SmtpPort)
	})
}

func TestConfigurationAttrValidate(t *testing.T) {
	t.Run("invalid verifier email", func(t *testing.T) {
		configurationAttr := ConfigurationAttr{VerifierEmail: "email@domain"}
		errorMessage := fmt.Sprintf("%v is invalid verifier email", configurationAttr.VerifierEmail)

		assert.EqualError(t, configurationAttr.validate(), errorMessage)
	})

	t.Run("invalid verifier domain", func(t *testing.T) {
		configurationAttr := ConfigurationAttr{
			VerifierEmail:  randomEmail(),
			VerifierDomain: "invalid_domain",
		}
		errorMessage := fmt.Sprintf("%v is invalid verifier domain", configurationAttr.VerifierDomain)

		assert.EqualError(t, configurationAttr.validate(), errorMessage)
	})

	t.Run("invalid default validation type", func(t *testing.T) {
		configurationAttr := ConfigurationAttr{VerifierEmail: randomEmail(), ValidationTypeDefault: "invalid validation type"}
		errorMessage := fmt.Sprintf("%v is invalid default validation type, use one of these: [regex mx mx_blacklist smtp]", configurationAttr.ValidationTypeDefault)

		assert.EqualError(t, configurationAttr.validate(), errorMessage)
	})

	t.Run("invalid connection timeout", func(t *testing.T) {
		configurationAttr := ConfigurationAttr{
			VerifierEmail:         randomEmail(),
			ValidationTypeDefault: randomValidationType(),
			ConnectionTimeout:     randomNegativeNumber(),
		}
		errorMessage := fmt.Sprintf("%v should be a positive integer", configurationAttr.ConnectionTimeout)

		assert.EqualError(t, configurationAttr.validate(), errorMessage)
	})

	t.Run("invalid response timeout", func(t *testing.T) {
		configurationAttr := ConfigurationAttr{
			VerifierEmail:         randomEmail(),
			ValidationTypeDefault: randomValidationType(),
			ConnectionTimeout:     randomPositiveNumber(),
			ResponseTimeout:       randomNegativeNumber(),
		}
		errorMessage := fmt.Sprintf("%v should be a positive integer", configurationAttr.ResponseTimeout)

		assert.EqualError(t, configurationAttr.validate(), errorMessage)
	})

	t.Run("invalid connection attempts", func(t *testing.T) {
		configurationAttr := ConfigurationAttr{
			VerifierEmail:         randomEmail(),
			ValidationTypeDefault: randomValidationType(),
			ConnectionTimeout:     randomPositiveNumber(),
			ResponseTimeout:       randomPositiveNumber(),
			ConnectionAttempts:    randomNegativeNumber(),
		}
		errorMessage := fmt.Sprintf("%v should be a positive integer", configurationAttr.ConnectionAttempts)

		assert.EqualError(t, configurationAttr.validate(), errorMessage)
	})

	t.Run("invalid SMTP port number", func(t *testing.T) {
		configurationAttr := ConfigurationAttr{
			VerifierEmail:         randomEmail(),
			ValidationTypeDefault: randomValidationType(),
			ConnectionTimeout:     randomPositiveNumber(),
			ResponseTimeout:       randomPositiveNumber(),
			ConnectionAttempts:    randomPositiveNumber(),
			SmtpPort:              randomNegativeNumber(),
		}
		errorMessage := fmt.Sprintf("%v should be a positive integer", configurationAttr.SmtpPort)

		assert.EqualError(t, configurationAttr.validate(), errorMessage)
	})

	t.Run("invalid whitelisted domains", func(t *testing.T) {
		configurationAttr := ConfigurationAttr{
			VerifierEmail:         randomEmail(),
			ValidationTypeDefault: randomValidationType(),
			ConnectionTimeout:     randomPositiveNumber(),
			ResponseTimeout:       randomPositiveNumber(),
			ConnectionAttempts:    randomPositiveNumber(),
			SmtpPort:              randomPositiveNumber(),
			WhitelistedDomains:    []string{randomDomain(), "a"},
		}
		errorMessage := fmt.Sprintf("%v is invalid domain name", configurationAttr.WhitelistedDomains[1])

		assert.EqualError(t, configurationAttr.validate(), errorMessage)
	})

	t.Run("invalid blacklisted domains", func(t *testing.T) {
		configurationAttr := ConfigurationAttr{
			VerifierEmail:         randomEmail(),
			ValidationTypeDefault: randomValidationType(),
			ConnectionTimeout:     randomPositiveNumber(),
			ResponseTimeout:       randomPositiveNumber(),
			ConnectionAttempts:    randomPositiveNumber(),
			SmtpPort:              randomPositiveNumber(),
			BlacklistedDomains:    []string{randomDomain(), "b"},
		}
		errorMessage := fmt.Sprintf("%v is invalid domain name", configurationAttr.BlacklistedDomains[1])

		assert.EqualError(t, configurationAttr.validate(), errorMessage)
	})

	t.Run("invalid blacklisted mx ip address", func(t *testing.T) {
		configurationAttr := ConfigurationAttr{
			VerifierEmail:            randomEmail(),
			ValidationTypeDefault:    randomValidationType(),
			ConnectionTimeout:        randomPositiveNumber(),
			ResponseTimeout:          randomPositiveNumber(),
			ConnectionAttempts:       randomPositiveNumber(),
			SmtpPort:                 randomPositiveNumber(),
			BlacklistedMxIpAddresses: []string{randomIpAddress(), "1.1.1.256:65536"},
		}
		errorMessage := fmt.Sprintf("%v is invalid ip address", configurationAttr.BlacklistedMxIpAddresses[1])

		assert.EqualError(t, configurationAttr.validate(), errorMessage)
	})

	t.Run("invalid dns, wrong ip address", func(t *testing.T) {
		configurationAttr := ConfigurationAttr{
			VerifierEmail:         randomEmail(),
			ValidationTypeDefault: randomValidationType(),
			ConnectionTimeout:     randomPositiveNumber(),
			ResponseTimeout:       randomPositiveNumber(),
			ConnectionAttempts:    randomPositiveNumber(),
			SmtpPort:              randomPositiveNumber(),
			Dns:                   "1.1.1.256",
		}
		errorMessage := fmt.Sprintf("%v is invalid dns server", configurationAttr.Dns)

		assert.EqualError(t, configurationAttr.validate(), errorMessage)
	})

	t.Run("invalid dns, wrong port number", func(t *testing.T) {
		configurationAttr := ConfigurationAttr{
			VerifierEmail:         randomEmail(),
			ValidationTypeDefault: randomValidationType(),
			ConnectionTimeout:     randomPositiveNumber(),
			ResponseTimeout:       randomPositiveNumber(),
			ConnectionAttempts:    randomPositiveNumber(),
			SmtpPort:              randomPositiveNumber(),
			Dns:                   "1.1.1.255:65536",
		}
		errorMessage := fmt.Sprintf("%v is invalid dns server", configurationAttr.Dns)

		assert.EqualError(t, configurationAttr.validate(), errorMessage)
	})

	t.Run("invalid dns, wrong ip address and port number", func(t *testing.T) {
		configurationAttr := ConfigurationAttr{
			VerifierEmail:         randomEmail(),
			ValidationTypeDefault: randomValidationType(),
			ConnectionTimeout:     randomPositiveNumber(),
			ResponseTimeout:       randomPositiveNumber(),
			ConnectionAttempts:    randomPositiveNumber(),
			SmtpPort:              randomPositiveNumber(),
			Dns:                   "1.1.1.256:65536",
		}
		errorMessage := fmt.Sprintf("%v is invalid dns server", configurationAttr.Dns)

		assert.EqualError(t, configurationAttr.validate(), errorMessage)
	})

	t.Run("invalid validation type by domain, wrong domain", func(t *testing.T) {
		invalidDomain := "inavlid domain"
		configurationAttr := ConfigurationAttr{
			VerifierEmail:          randomEmail(),
			ValidationTypeDefault:  randomValidationType(),
			ConnectionTimeout:      randomPositiveNumber(),
			ResponseTimeout:        randomPositiveNumber(),
			ConnectionAttempts:     randomPositiveNumber(),
			SmtpPort:               randomPositiveNumber(),
			ValidationTypeByDomain: map[string]string{randomDomain(): "regex", invalidDomain: "wrong_type"},
		}
		errorMessage := fmt.Sprintf("%v is invalid domain name", invalidDomain)

		assert.EqualError(t, configurationAttr.validate(), errorMessage)
	})

	t.Run("invalid validation type by domain, wrong validation type", func(t *testing.T) {
		invalidType := "inavlid validation type"
		configurationAttr := ConfigurationAttr{
			VerifierEmail:          randomEmail(),
			ValidationTypeDefault:  randomValidationType(),
			ConnectionTimeout:      randomPositiveNumber(),
			ResponseTimeout:        randomPositiveNumber(),
			ConnectionAttempts:     randomPositiveNumber(),
			SmtpPort:               randomPositiveNumber(),
			ValidationTypeByDomain: map[string]string{randomDomain(): "regex", randomDomain(): invalidType},
		}
		errorMessage := fmt.Sprintf("%v is invalid default validation type, use one of these: [regex mx mx_blacklist smtp]", invalidType)

		assert.EqualError(t, configurationAttr.validate(), errorMessage)
	})

	t.Run("invalid email pattern", func(t *testing.T) {
		configurationAttr := ConfigurationAttr{
			VerifierEmail:         randomEmail(),
			ValidationTypeDefault: randomValidationType(),
			ConnectionTimeout:     randomPositiveNumber(),
			ResponseTimeout:       randomPositiveNumber(),
			ConnectionAttempts:    randomPositiveNumber(),
			SmtpPort:              randomPositiveNumber(),
			EmailPattern:          `\K`,
		}
		errorMessage := fmt.Sprintf("error parsing regexp: invalid escape sequence: `%v`", configurationAttr.EmailPattern)

		assert.EqualError(t, configurationAttr.validate(), errorMessage)
	})

	t.Run("invalid smtp error body pattern", func(t *testing.T) {
		configurationAttr := ConfigurationAttr{
			VerifierEmail:         randomEmail(),
			ValidationTypeDefault: randomValidationType(),
			ConnectionTimeout:     randomPositiveNumber(),
			ResponseTimeout:       randomPositiveNumber(),
			ConnectionAttempts:    randomPositiveNumber(),
			SmtpPort:              randomPositiveNumber(),
			SmtpErrorBodyPattern:  `\K`,
		}
		errorMessage := fmt.Sprintf("error parsing regexp: invalid escape sequence: `%v`", configurationAttr.SmtpErrorBodyPattern)

		assert.EqualError(t, configurationAttr.validate(), errorMessage)
	})

	t.Run("coerces to special format/types", func(t *testing.T) {
		randomIpAddress, randomIpAddressWithDefaultPortNumber := randomDnsServerWithDefaultPortNumber()
		regexPatternFirst, regexPatternSecond := "1", "2"
		configurationAttr := ConfigurationAttr{
			VerifierEmail:         randomEmail(),
			ValidationTypeDefault: randomValidationType(),
			ConnectionTimeout:     randomPositiveNumber(),
			ResponseTimeout:       randomPositiveNumber(),
			ConnectionAttempts:    randomPositiveNumber(),
			SmtpPort:              randomPositiveNumber(),
			Dns:                   randomIpAddress,
			EmailPattern:          regexPatternFirst,
			SmtpErrorBodyPattern:  regexPatternSecond,
		}
		regexObjectFirst, _ := newRegex(regexPatternFirst)
		regexObjectSecond, _ := newRegex(regexPatternSecond)

		assert.NoError(t, configurationAttr.validate())
		assert.Equal(t, randomIpAddressWithDefaultPortNumber, configurationAttr.Dns)
		assert.EqualValues(t, regexObjectFirst, configurationAttr.RegexEmail)
		assert.EqualValues(t, regexObjectSecond, configurationAttr.RegexSmtpErrorBody)
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
		actualDomain, err := new(ConfigurationAttr).buildVerifierDomain(verifierEmail, emptyString)

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
		assert.NoError(t, new(ConfigurationAttr).validateDnsServerContext(randomIpAddress()))
	})

	t.Run("valid dns server ip with port number", func(t *testing.T) {
		assert.NoError(t, new(ConfigurationAttr).validateDnsServerContext(randomIpAddress()+":65507"))
	})

	t.Run("invalid dns server ip without port number", func(t *testing.T) {
		invalidDNSServer := "1.1.1.256"
		errorMessage := fmt.Sprintf("%s is invalid dns server", invalidDNSServer)

		assert.EqualError(t, new(ConfigurationAttr).validateDnsServerContext(invalidDNSServer), errorMessage)
	})

	t.Run("valid dns server ip with invalid port number", func(t *testing.T) {
		invalidDNSServer := "1.1.1.1:65536"
		errorMessage := fmt.Sprintf("%s is invalid dns server", invalidDNSServer)

		assert.EqualError(t, new(ConfigurationAttr).validateDnsServerContext(invalidDNSServer), errorMessage)
	})

	t.Run("invalid dns server ip with invalid port number", func(t *testing.T) {
		invalidDNSServer := "256.256.256.256:0"
		errorMessage := fmt.Sprintf("%s is invalid dns server", invalidDNSServer)

		assert.EqualError(t, new(ConfigurationAttr).validateDnsServerContext(invalidDNSServer), errorMessage)
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

func TestConfigurationAttrValidateWithFormatDnsServerContext(t *testing.T) {
	t.Run("when DNS gateway not specified", func(t *testing.T) {
		dnsGateway, err := new(ConfigurationAttr).validateWithFormatDnsServerContext(emptyString)

		assert.Equal(t, emptyString, dnsGateway)
		assert.NoError(t, err)
	})

	t.Run("when DNS gateway specified without port", func(t *testing.T) {
		dns := randomIpAddress()
		dnsGateway, err := new(ConfigurationAttr).validateWithFormatDnsServerContext(dns)

		assert.Equal(t, serverWithPortNumber(dns, defaultDnsPort), dnsGateway)
		assert.NoError(t, err)
	})

	t.Run("when DNS gateway specified with port", func(t *testing.T) {
		dns := serverWithPortNumber(randomIpAddress(), randomPortNumber())
		dnsGateway, err := new(ConfigurationAttr).validateWithFormatDnsServerContext(dns)

		assert.Equal(t, dns, dnsGateway)
		assert.NoError(t, err)
	})

	t.Run("when DNS gateway specified, wrong ip address", func(t *testing.T) {
		dns := "1.1.1.256"
		errorMessage := fmt.Sprintf("%v is invalid dns server", dns)
		dnsGateway, err := new(ConfigurationAttr).validateWithFormatDnsServerContext(dns)

		assert.Equal(t, dns, dnsGateway)
		assert.EqualError(t, err, errorMessage)
	})

	t.Run("when DNS gateway specified, wrong port number", func(t *testing.T) {
		dns := "1.1.1.255:65536"
		errorMessage := fmt.Sprintf("%v is invalid dns server", dns)
		dnsGateway, err := new(ConfigurationAttr).validateWithFormatDnsServerContext(dns)

		assert.Equal(t, dns, dnsGateway)
		assert.EqualError(t, err, errorMessage)
	})

	t.Run("when DNS gateway specified, wrong ip address and port number", func(t *testing.T) {
		dns := "1.1.1.256:65536"
		errorMessage := fmt.Sprintf("%v is invalid dns server", dns)
		dnsGateway, err := new(ConfigurationAttr).validateWithFormatDnsServerContext(dns)

		assert.Equal(t, dns, dnsGateway)
		assert.EqualError(t, err, errorMessage)
	})
}
