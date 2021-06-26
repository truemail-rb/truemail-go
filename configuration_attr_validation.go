package truemail

import (
	"fmt"
	"regexp"
)

func validateVerifierEmail(verifierEmail string) error {
	if matchRegex(verifierEmail, EmailCharsSize) && matchRegex(verifierEmail, RegexEmailPattern) {
		return nil
	}
	return fmt.Errorf("%s is invalid verifier email", verifierEmail)
}

func validateVerifierDomain(verifierDomain string) (string, error) {
	if matchRegex(verifierDomain, DomainCharsSize) && matchRegex(verifierDomain, RegexDomainPattern) {
		return verifierDomain, nil
	}
	return verifierDomain, fmt.Errorf("%s is invalid verifier domain", verifierDomain)
}

func buildVerifierDomain(verifierEmail, verifierDomain string) (string, error) {
	if verifierDomain == "" {
		regex, _ := newRegex(RegexEmailPattern)
		domainCaptureGroup := 3
		return regex.FindStringSubmatch(verifierEmail)[domainCaptureGroup], nil
	}
	return validateVerifierDomain(verifierDomain)
}

func availableValidationTypes() []string {
	return []string{ValidationTypeRegex, ValidationTypeMx, ValidationTypeSMTP}
}

func validateValidationTypeDefaultContext(ValidationTypeDefault string) error {
	if isIncluded(availableValidationTypes(), ValidationTypeDefault) {
		return nil
	}
	return fmt.Errorf(
		"%s is invalid default validation type, use one of these: %s",
		ValidationTypeDefault,
		availableValidationTypes(),
	)
}

func validateIntegerPositive(integer int) error {
	if integer > 0 {
		return nil
	}
	return fmt.Errorf("%v should be a positive integer", integer)
}

func validateStringContext(target, regexPattern, msg string) error {
	if matchRegex(target, regexPattern) {
		return nil
	}
	return fmt.Errorf("%s is invalid %s", target, msg)
}

func validateDomainContext(domainName string) error {
	return validateStringContext(domainName, RegexDomainPattern, "domain name")
}

func validateDomainsContext(domains []string) error {
	for _, domainName := range domains {
		err := validateDomainContext(domainName)
		if err != nil {
			return err
		}
	}
	return nil
}

func validateIpAddressContext(ipAddress string) error {
	return validateStringContext(ipAddress, RegexIpAddressPattern, "ip address")
}

func validateIpAddressesContext(ipAddresses []string) error {
	for _, ipAddress := range ipAddresses {
		err := validateIpAddressContext(ipAddress)
		if err != nil {
			return err
		}
	}
	return nil
}

func validateDNSServerContext(dnsServer string) error {
	return validateStringContext(dnsServer, RegexDNSServerAddressPattern, "dns server")
}

func validateDNSServersContext(dnsServers []string) error {
	for _, dnsServer := range dnsServers {
		err := validateDNSServerContext(dnsServer)
		if err != nil {
			return err
		}
	}
	return nil
}

func validateTypeByDomainContext(typesByDomains map[string]string) error {
	for domainName, validationType := range typesByDomains {
		err := validateDomainContext(domainName)
		if err != nil {
			return err
		}

		err = validateValidationTypeDefaultContext(validationType)
		if err != nil {
			return err
		}
	}
	return nil
}

func isIncluded(slice []string, target string) bool {
	for _, item := range slice {
		if item == target {
			return true
		}
	}
	return false
}

func newRegex(regexPattern string) (*regexp.Regexp, error) {
	return regexp.Compile(regexPattern)
}

func matchRegex(strContext, regexPattern string) bool {
	regex, err := newRegex(regexPattern)
	if err != nil {
		return false
	}
	return regex.MatchString(strContext)
}
