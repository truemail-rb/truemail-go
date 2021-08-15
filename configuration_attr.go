package truemail

import (
	"context"
	"fmt"
	"regexp"
)

// ConfigurationAttr kwargs structure for configuration builder
type ConfigurationAttr struct {
	ctx                                                                                           context.Context
	verifierEmail, verifierDomain, validationTypeDefault, emailPattern, smtpErrorBodyPattern, dns string
	connectionTimeout, responseTimeout, connectionAttempts                                        int
	whitelistedDomains, blacklistedDomains, blacklistedMxIpAddresses                              []string
	validationTypeByDomain                                                                        map[string]string
	whitelistValidation, notRfcMxLookupFlow, smtpFailFast, smtpSafeCheck                          bool
	regexEmail, regexSMTPErrorBody                                                                *regexp.Regexp
}

// ConfigurationAttr methods

// assigns default values to ConfigurationAttr fields
func (config *ConfigurationAttr) assignDefaultValues() {
	if config.validationTypeDefault == EmptyString {
		config.validationTypeDefault = ValidationTypeDefault
	}
	if config.emailPattern == EmptyString {
		config.emailPattern = RegexEmailPattern
	}
	if config.smtpErrorBodyPattern == EmptyString {
		config.smtpErrorBodyPattern = RegexSMTPErrorBodyPattern
	}
	if config.connectionTimeout == 0 {
		config.connectionTimeout = DefaultConnectionTimeout
	}
	if config.responseTimeout == 0 {
		config.responseTimeout = DefaultResponseTimeout
	}
	if config.connectionAttempts == 0 {
		config.connectionAttempts = DefaultConnectionAttempts
	}
}

// validates and coerces ConfigurationAttr fields context
func (config *ConfigurationAttr) validate() error {
	err := config.validateVerifierEmail(config.verifierEmail)
	if err != nil {
		return err
	}

	config.verifierDomain, err = config.buildVerifierDomain(config.verifierEmail, config.verifierDomain)
	if err != nil {
		return err
	}

	err = config.validateValidationTypeDefaultContext(config.validationTypeDefault)
	if err != nil {
		return err
	}

	err = config.validateIntegerPositive(config.connectionTimeout)
	if err != nil {
		return err
	}

	err = config.validateIntegerPositive(config.responseTimeout)
	if err != nil {
		return err
	}

	err = config.validateIntegerPositive(config.connectionAttempts)
	if err != nil {
		return err
	}

	err = config.validateDomainsContext(config.whitelistedDomains)
	if err != nil {
		return err
	}

	err = config.validateDomainsContext(config.blacklistedDomains)
	if err != nil {
		return err
	}

	err = config.validateIpAddressesContext(config.blacklistedMxIpAddresses)
	if err != nil {
		return err
	}

	dns := config.dns

	if dns != EmptyString {
		err = config.validateDNSServerContext(dns)
		if err != nil {
			return err
		}

		config.dns = config.formatDns(dns)
	}

	err = config.validateTypeByDomainContext(config.validationTypeByDomain)
	if err != nil {
		return err
	}

	config.regexEmail, err = newRegex(config.emailPattern)
	if err != nil {
		return err
	}

	config.regexSMTPErrorBody, err = newRegex(config.smtpErrorBodyPattern)
	if err != nil {
		return err
	}

	return nil
}

// Validates verifier email. Returns error if validation fails
func (config *ConfigurationAttr) validateVerifierEmail(verifierEmail string) error {
	if matchRegex(verifierEmail, EmailCharsSize) && matchRegex(verifierEmail, RegexEmailPattern) {
		return nil
	}
	return fmt.Errorf("%s is invalid verifier email", verifierEmail)
}

// Validates verifier domain. Returns error if validation fails
func (config *ConfigurationAttr) validateVerifierDomain(verifierDomain string) (string, error) {
	if matchRegex(verifierDomain, DomainCharsSize) && matchRegex(verifierDomain, RegexDomainPattern) {
		return verifierDomain, nil
	}
	return verifierDomain, fmt.Errorf("%s is invalid verifier domain", verifierDomain)
}

// Returns verifier domain builded from verifier email or validated from function arg
func (config *ConfigurationAttr) buildVerifierDomain(verifierEmail, verifierDomain string) (string, error) {
	if verifierDomain == EmptyString {
		regex, _ := newRegex(RegexEmailPattern)
		domainCaptureGroup := 3
		return regex.FindStringSubmatch(verifierEmail)[domainCaptureGroup], nil
	}
	return config.validateVerifierDomain(verifierDomain)
}

// Validates validation type. Returns error if validation fails
func (config *ConfigurationAttr) validateValidationTypeDefaultContext(ValidationTypeDefault string) error {
	if isIncluded(availableValidationTypes(), ValidationTypeDefault) {
		return nil
	}
	return fmt.Errorf(
		"%s is invalid default validation type, use one of these: %s",
		ValidationTypeDefault,
		availableValidationTypes(),
	)
}

// Validates is integer is a positive. Returns error if validation fails
func (config *ConfigurationAttr) validateIntegerPositive(integer int) error {
	if integer > 0 {
		return nil
	}
	return fmt.Errorf("%v should be a positive integer", integer)
}

// Validates is string matches to regex pattern. Returns error if validation fails
func (config *ConfigurationAttr) validateStringContext(target, regexPattern, msg string) error {
	if matchRegex(target, regexPattern) {
		return nil
	}
	return fmt.Errorf("%s is invalid %s", target, msg)
}

// Validates is domain name matches to regex domain pattern.
// Returns error if validation fails
func (config *ConfigurationAttr) validateDomainContext(domainName string) error {
	return config.validateStringContext(domainName, RegexDomainPattern, "domain name")
}

// Validates is each domain name from slice matches to regex domain pattern.
// Returns error if at least one of domain validations fails
func (config *ConfigurationAttr) validateDomainsContext(domains []string) error {
	for _, domainName := range domains {
		err := config.validateDomainContext(domainName)
		if err != nil {
			return err
		}
	}
	return nil
}

// Validates is ip address matches to regex ip address pattern.
// Returns error if validation fails
func (config *ConfigurationAttr) validateIpAddressContext(ipAddress string) error {
	return config.validateStringContext(ipAddress, RegexIpAddressPattern, "ip address")
}

// Validates is ip address matches to regex ip address pattern.
// Returns error if at least one of ip address validations fails
func (config *ConfigurationAttr) validateIpAddressesContext(ipAddresses []string) error {
	for _, ipAddress := range ipAddresses {
		err := config.validateIpAddressContext(ipAddress)
		if err != nil {
			return err
		}
	}
	return nil
}

// Validates is DNS server matches to regex DNS server address pattern.
// Returns error if validation fails
func (config *ConfigurationAttr) validateDNSServerContext(dnsServer string) error {
	return config.validateStringContext(dnsServer, RegexDNSServerAddressPattern, "dns server")
}

// Validates typesByDomains map key-values. Returns error if validation fails
func (config *ConfigurationAttr) validateTypeByDomainContext(typesByDomains map[string]string) error {
	for domainName, validationType := range typesByDomains {
		err := config.validateDomainContext(domainName)
		if err != nil {
			return err
		}

		err = config.validateValidationTypeDefaultContext(validationType)
		if err != nil {
			return err
		}
	}
	return nil
}

// Addes default DNS port to ip address by template {ipAddress}:{portNumber} for cases
// when port number is not specified
func (config *ConfigurationAttr) formatDns(dnsGateway string) string {
	regex, _ := newRegex(RegexDNSServerAddressPattern)
	portNumberCaptureGroup := 5
	if regex.FindStringSubmatch(dnsGateway)[portNumberCaptureGroup] != EmptyString {
		return dnsGateway
	}

	return dnsGateway + ":" + DefaultDnsPort
}
