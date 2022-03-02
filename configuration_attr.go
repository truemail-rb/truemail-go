package truemail

import (
	"context"
	"fmt"
	"regexp"
)

// ConfigurationAttr kwargs structure for configuration builder
type ConfigurationAttr struct {
	ctx                                                                                           context.Context
	VerifierEmail, VerifierDomain, ValidationTypeDefault, EmailPattern, SmtpErrorBodyPattern, Dns string
	ConnectionTimeout, ResponseTimeout, ConnectionAttempts, SmtpPort                              int
	WhitelistedDomains, BlacklistedDomains, BlacklistedMxIpAddresses                              []string
	ValidationTypeByDomain                                                                        map[string]string
	WhitelistValidation, NotRfcMxLookupFlow, SmtpFailFast, SmtpSafeCheck                          bool
	RegexEmail, RegexSmtpErrorBody                                                                *regexp.Regexp
}

// ConfigurationAttr methods

// assigns default values to ConfigurationAttr fields
func (config *ConfigurationAttr) assignDefaultValues() {
	if config.ValidationTypeDefault == emptyString {
		config.ValidationTypeDefault = validationTypeDefault
	}
	if config.EmailPattern == emptyString {
		config.EmailPattern = regexEmailPattern
	}
	if config.SmtpErrorBodyPattern == emptyString {
		config.SmtpErrorBodyPattern = regexSMTPErrorBodyPattern
	}
	if config.ConnectionTimeout == 0 {
		config.ConnectionTimeout = defaultConnectionTimeout
	}
	if config.ResponseTimeout == 0 {
		config.ResponseTimeout = defaultResponseTimeout
	}
	if config.ConnectionAttempts == 0 {
		config.ConnectionAttempts = defaultConnectionAttempts
	}

	if config.SmtpPort == 0 {
		config.SmtpPort = defaultSmtpPort
	}
}

// validates and coerces ConfigurationAttr fields context
func (config *ConfigurationAttr) validate() error {
	err := config.validateVerifierEmail(config.VerifierEmail)
	if err != nil {
		return err
	}

	config.VerifierDomain, err = config.buildVerifierDomain(config.VerifierEmail, config.VerifierDomain)
	if err != nil {
		return err
	}

	err = config.validateValidationTypeDefaultContext(config.ValidationTypeDefault)
	if err != nil {
		return err
	}

	err = config.validateIntegerPositive(config.ConnectionTimeout)
	if err != nil {
		return err
	}

	err = config.validateIntegerPositive(config.ResponseTimeout)
	if err != nil {
		return err
	}

	err = config.validateIntegerPositive(config.ConnectionAttempts)
	if err != nil {
		return err
	}

	err = config.validateIntegerPositive(config.SmtpPort)
	if err != nil {
		return err
	}

	err = config.validateDomainsContext(config.WhitelistedDomains)
	if err != nil {
		return err
	}

	err = config.validateDomainsContext(config.BlacklistedDomains)
	if err != nil {
		return err
	}

	err = config.validateIpAddressesContext(config.BlacklistedMxIpAddresses)
	if err != nil {
		return err
	}

	dns, err := config.validateWithFormatDnsServerContext(config.Dns)
	if err != nil {
		return err
	}

	config.Dns = dns

	err = config.validateTypeByDomainContext(config.ValidationTypeByDomain)
	if err != nil {
		return err
	}

	config.RegexEmail, err = newRegex(config.EmailPattern)
	if err != nil {
		return err
	}

	config.RegexSmtpErrorBody, err = newRegex(config.SmtpErrorBodyPattern)
	if err != nil {
		return err
	}

	return nil
}

// Validates verifier email. Returns error if validation fails
func (config *ConfigurationAttr) validateVerifierEmail(verifierEmail string) error {
	if matchRegex(verifierEmail, emailCharsSize) && matchRegex(verifierEmail, regexEmailPattern) {
		return nil
	}
	return fmt.Errorf("%s is invalid verifier email", verifierEmail)
}

// Validates verifier domain. Returns error if validation fails
func (config *ConfigurationAttr) validateVerifierDomain(verifierDomain string) (string, error) {
	if matchRegex(verifierDomain, domainCharsSize) && matchRegex(verifierDomain, regexDomainPattern) {
		return verifierDomain, nil
	}
	return verifierDomain, fmt.Errorf("%s is invalid verifier domain", verifierDomain)
}

// Returns verifier domain builded from verifier email or validated from function arg
func (config *ConfigurationAttr) buildVerifierDomain(verifierEmail, verifierDomain string) (string, error) {
	if verifierDomain == emptyString {
		regex, _ := newRegex(regexEmailPattern)
		domainCaptureGroup := 3
		return regex.FindStringSubmatch(verifierEmail)[domainCaptureGroup], nil
	}
	return config.validateVerifierDomain(verifierDomain)
}

// Validates validation type. Returns error if validation fails
func (config *ConfigurationAttr) validateValidationTypeDefaultContext(validationTypeDefault string) error {
	if isIncluded(availableValidationTypes(), validationTypeDefault) {
		return nil
	}
	return fmt.Errorf(
		"%s is invalid default validation type, use one of these: %s",
		validationTypeDefault,
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
	return config.validateStringContext(domainName, regexDomainPattern, "domain name")
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
	return config.validateStringContext(ipAddress, regexIpAddressPattern, "ip address")
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
func (config *ConfigurationAttr) validateDnsServerContext(dnsServer string) error {
	return config.validateStringContext(dnsServer, regexDNSServerAddressPattern, "dns server")
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
	regex, _ := newRegex(regexDNSServerAddressPattern)
	portNumberCaptureGroup := 5
	if regex.FindStringSubmatch(dnsGateway)[portNumberCaptureGroup] != emptyString {
		return dnsGateway
	}

	return serverWithPortNumber(dnsGateway, defaultDnsPort)
}

// Validates DNS server context and returns formatted DNS when dns gateway not empty.
// Otherwise returns empty string
func (config *ConfigurationAttr) validateWithFormatDnsServerContext(dnsGateway string) (string, error) {
	if dnsGateway == emptyString {
		return emptyString, nil
	}

	err := config.validateDnsServerContext(dnsGateway)
	if err != nil {
		return dnsGateway, err
	}

	return config.formatDns(dnsGateway), nil
}
