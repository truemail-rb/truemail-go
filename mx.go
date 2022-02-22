package truemail

import (
	"fmt"

	"golang.org/x/net/idna"
)

// DNS (MX) validation resolver interface
type resolver interface {
	aRecord(string) (string, error)
	aRecords(string) ([]string, error)
	cnameRecord(string) (string, error)
	ptrRecords(string) ([]string, error)
	mxRecords(string) ([]uint16, []string, error)
}

// DNS (MX) validation, second validation level
type validationMx struct {
	result *validatorResult
	resolver
}

// interface implementation
func (validation *validationMx) check(validatorResult *validatorResult) *validatorResult {
	validation.result = validatorResult
	validation.setValidatorResultPunycodeRepresentation()
	validation.initDnsResolver()
	validation.runMxLookup()

	if validation.isMailServerNotFound() {
		validatorResult.Success = false
		validatorResult.addError(validationTypeMx, mxErrorContext)
	}

	return validatorResult
}

// validationMx methods

// Returns punycode domain representation
func (validation *validationMx) punycodeDomain(domain string) string {
	punycodeDomain, _ := idna.New().ToASCII(domain)
	return punycodeDomain
}

// Assigns punycodeEmail, punycodeDomain representations to validatorResult
func (validation *validationMx) setValidatorResultPunycodeRepresentation() {
	email := validation.result.Email
	user := regexCaptureGroup(email, regexEmailPattern, 2)
	punycodeDomain := validation.punycodeDomain(regexCaptureGroup(email, regexEmailPattern, 3))

	validation.result.punycodeEmail = user + "@" + punycodeDomain
	validation.result.punycodeDomain = punycodeDomain
}

// Initializes MX validation DNS resolver
func (validation *validationMx) initDnsResolver() {
	validation.resolver = newDnsResolver(validation.result.Configuration)
}

// Returns true if validatorResult contains no mail servers, otherwise returns false
func (validation *validationMx) isMailServerNotFound() bool {
	return len(validation.result.MailServers) == 0
}

// Returns true if validatorResult contains mail servers, otherwise returns false
func (validation *validationMx) isMailServerFound() bool {
	return len(validation.result.MailServers) > 0
}

// Addes just uniques hosts to validatorResult.MailServers
func (validation *validationMx) fetchTargetHosts(hosts ...string) {
	mailServers := validation.result.MailServers
	validation.result.MailServers = append(mailServers, sliceDiff(uniqStrings(hosts), mailServers)...)
}

// Returns true if connection attempts more than zero, otherwise returns false
func (validation *validationMx) isConnectionAttemptsAvailable(connectionAttempts int) bool {
	return connectionAttempts > 0
}

// Casts is wrapped error is an DnsNotFound error
func (validation *validationMx) isDnsNotFoundError(err error) bool {
	e, ok := err.(*validationError)
	return ok && e.isDnsNotFound
}

// Casts is wrapped error is an NullMxError error
func (validation *validationMx) isNullMxError(err error) bool {
	e, ok := err.(*validationError)
	return ok && e.isNullMxFound
}

// A records resolver, the part of MX records resolver
func (validation *validationMx) aRecords(hostName string) (ipAddresses []string, err error) {
	connectionAttempts := validation.result.Configuration.ConnectionAttempts
	for validation.isConnectionAttemptsAvailable(connectionAttempts) {
		ipAddresses, err = validation.resolver.aRecords(hostName)
		connectionAttempts -= 1

		if err == nil {
			break
		} else {
			if validation.isDnsNotFoundError(err) {
				break
			}
		}
	}

	return ipAddresses, err
}

// MX records resolver
func (validation *validationMx) hostsFromMxRecords(hostName string) (resolvedIpAddresses []string, err error) {
	var priorities []uint16
	var hostNames, ipAddresses []string
	connectionAttempts := validation.result.Configuration.ConnectionAttempts

	// Resolves MX hostnames by hostname
	for validation.isConnectionAttemptsAvailable(connectionAttempts) {
		priorities, hostNames, err = validation.resolver.mxRecords(hostName)
		connectionAttempts -= 1

		if err == nil {
			break
		} else {
			if validation.isDnsNotFoundError(err) {
				return resolvedIpAddresses, err
			}
		}
	}

	// Checkes null MX record
	if len(hostNames) == 1 && priorities[0] == 0 && hostNames[0] == emptyString {
		return resolvedIpAddresses, wrapNullMxError(fmt.Errorf("%s includes null MX record", hostName))
	}

	// Resolves host addresses by MX hostname
	for _, hostName := range hostNames {
		ipAddresses, err = validation.aRecords(hostName)
		if err != nil {
			continue
		}

		resolvedIpAddresses = append(resolvedIpAddresses, ipAddresses...)
	}

	return resolvedIpAddresses, err
}

// A record resolver
func (validation *validationMx) hostFromARecord(hostName string) (resolvedIpAddress string, err error) {
	connectionAttempts := validation.result.Configuration.ConnectionAttempts
	for validation.isConnectionAttemptsAvailable(connectionAttempts) {
		resolvedIpAddress, err = validation.resolver.aRecord(hostName)
		connectionAttempts -= 1

		if err == nil {
			break
		} else {
			if validation.isDnsNotFoundError(err) {
				break
			}
		}
	}

	return resolvedIpAddress, err
}

func (validation *validationMx) ptrRecords(hostAddress string) (resolvedHostNames []string, err error) {
	connectionAttempts := validation.result.Configuration.ConnectionAttempts
	for validation.isConnectionAttemptsAvailable(connectionAttempts) {
		resolvedHostNames, err = validation.resolver.ptrRecords(hostAddress)
		connectionAttempts -= 1

		if err == nil {
			break
		} else {
			if validation.isDnsNotFoundError(err) {
				break
			}
		}
	}

	return resolvedHostNames, err
}

// CNAME record resolver
func (validation *validationMx) hostsFromCnameRecord(hostName string) (resolvedIpAddresses []string, err error) {
	var resolvedHostNameByCname, resolvedIpAddressByARecord string
	var resolvedHostNamesByPtrRecords, resolvedIpAddressesByMxRecords []string
	connectionAttempts := validation.result.Configuration.ConnectionAttempts

	// Resolves hostname by CNAME record
	for validation.isConnectionAttemptsAvailable(connectionAttempts) {
		resolvedHostNameByCname, err = validation.resolver.cnameRecord(hostName)
		connectionAttempts -= 1

		if err == nil {
			break
		} else {
			if validation.isDnsNotFoundError(err) {
				break
			}
		}
	}

	if err != nil {
		return resolvedIpAddresses, err
	}

	// Resolves host address by A record
	resolvedIpAddressByARecord, err = validation.hostFromARecord(resolvedHostNameByCname)
	if err != nil {
		return resolvedIpAddresses, err
	}

	// Resolves hostnames by PTR records
	resolvedHostNamesByPtrRecords, err = validation.ptrRecords(resolvedIpAddressByARecord)
	if err != nil {
		return resolvedIpAddresses, err
	}

	// Resolves host addresses by MX records
	for _, resolvedHostName := range resolvedHostNamesByPtrRecords {
		resolvedIpAddressesByMxRecords, err = validation.hostsFromMxRecords(resolvedHostName)
		if err != nil {
			continue
		}

		resolvedIpAddresses = append(resolvedIpAddresses, resolvedIpAddressesByMxRecords...)
	}

	return resolvedIpAddresses, err
}

// Complex MX lookup for target domain, uses step by step MX, CNAME and A resolvers
func (validation *validationMx) runMxLookup() {
	var hostAddress string
	var hostAddresses []string
	var err error
	targetHostname := validation.result.punycodeDomain

	// MX records resolver
	hostAddresses, err = validation.hostsFromMxRecords(targetHostname)
	if err == nil {
		validation.fetchTargetHosts(hostAddresses...)
		return
	}

	if validation.isNullMxError(err) || validation.result.Configuration.NotRfcMxLookupFlow {
		return
	}

	// CNAME record resolver
	hostAddresses, err = validation.hostsFromCnameRecord(targetHostname)
	if err == nil {
		validation.fetchTargetHosts(hostAddresses...)
		return
	}

	// A record resolver
	hostAddress, err = validation.hostFromARecord(targetHostname)
	if err == nil {
		validation.fetchTargetHosts(hostAddress)
		return
	}
}
