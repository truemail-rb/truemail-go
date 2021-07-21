package truemail

import (
	"context"
	"regexp"
)

// configuration structure
type configuration struct {
	ctx                                                                  context.Context
	VerifierEmail, VerifierDomain, ValidationTypeDefault, DNS            string
	ConnectionTimeout, ResponseTimeout, ConnectionAttempts               int
	WhitelistedDomains, BlacklistedDomains, BlacklistedMxIpAddresses     []string
	ValidationTypeByDomain                                               map[string]string
	WhitelistValidation, NotRfcMxLookupFlow, SMTPFailFast, SMTPSafeCheck bool
	EmailPattern, SMTPErrorBodyPattern                                   *regexp.Regexp
}

// ConfigurationAttr kwargs for configuration builder
type ConfigurationAttr struct {
	ctx                                                                                           context.Context
	verifierEmail, verifierDomain, validationTypeDefault, emailPattern, smtpErrorBodyPattern, dns string
	connectionTimeout, responseTimeout, connectionAttempts                                        int
	whitelistedDomains, blacklistedDomains, blacklistedMxIpAddresses                              []string
	validationTypeByDomain                                                                        map[string]string
	whitelistValidation, notRfcMxLookupFlow, smtpFailFast, smtpSafeCheck                          bool
}

// Configuration builder constants, regex patterns
const (
	DefaultConnectionTimeout      = 2
	DefaultResponseTimeout        = 2
	DefaultConnectionAttempts     = 2
	ValidationTypeDomainListMatch = "domain_list_match"
	ValidationTypeRegex           = "regex"
	ValidationTypeMx              = "mx"
	ValidationTypeMxBlacklist     = "mx_blacklist"
	ValidationTypeSMTP            = "smtp"
	ValidationTypeDefault         = ValidationTypeSMTP
	DomainCharsSize               = `\A.{4,255}\z`
	EmailCharsSize                = `\A.{6,255}\z`
	RegexDomainPattern            = `(?i)[\p{L}0-9]+([\-.]{1}[\p{L}0-9]+)*\.\p{L}{2,63}`
	RegexEmailPattern             = `(\A([\p{L}0-9]+[\W\w]*)@(` + RegexDomainPattern + `)\z)`
	RegexDomainFromEmail          = `\A.+@(.+)\z`
	RegexSMTPErrorBodyPattern     = `(?i).*550{1}.*(user|account|customer|mailbox).*`
	RegexPortNumber               = `(6553[0-5]|655[0-2]\d|65[0-4](\d){2}|6[0-4](\d){3}|[1-5](\d){4}|[1-9](\d){0,3})`
	RegexIpAddress                = `((\d|[1-9]\d|1\d{2}|2[0-4]\d|25[0-5])\.){3}(\d|[1-9]\d|1\d{2}|2[0-4]\d|25[0-5])`
	RegexIpAddressPattern         = `\A` + RegexIpAddress + `\z`
	RegexDNSServerAddressPattern  = `\A` + RegexIpAddress + `(:` + RegexPortNumber + `)?\z`
	EmptyString                   = ""
)

// NewConfiguration builder
func NewConfiguration(config ConfigurationAttr) (*configuration, error) {
	// assign fileds default values
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

	// validate fileds context
	err := validateVerifierEmail(config.verifierEmail)
	if err != nil {
		return nil, err
	}

	config.verifierDomain, err = buildVerifierDomain(config.verifierEmail, config.verifierDomain)
	if err != nil {
		return nil, err
	}

	err = validateValidationTypeDefaultContext(config.validationTypeDefault)
	if err != nil {
		return nil, err
	}

	err = validateIntegerPositive(config.connectionTimeout)
	if err != nil {
		return nil, err
	}

	err = validateIntegerPositive(config.responseTimeout)
	if err != nil {
		return nil, err
	}

	err = validateIntegerPositive(config.connectionAttempts)
	if err != nil {
		return nil, err
	}

	err = validateDomainsContext(config.whitelistedDomains)
	if err != nil {
		return nil, err
	}

	err = validateDomainsContext(config.blacklistedDomains)
	if err != nil {
		return nil, err
	}

	err = validateIpAddressesContext(config.blacklistedMxIpAddresses)
	if err != nil {
		return nil, err
	}

	dns := config.dns

	if dns != EmptyString {
		err = validateDNSServerContext(config.dns)
		if err != nil {
			return nil, err
		}
	}

	err = validateTypeByDomainContext(config.validationTypeByDomain)
	if err != nil {
		return nil, err
	}

	regexEmail, err := newRegex(config.emailPattern)
	if err != nil {
		return nil, err
	}

	regexSMTPErrorBody, err := newRegex(config.smtpErrorBodyPattern)
	if err != nil {
		return nil, err
	}

	// create new configuration
	newConfiguration := configuration{
		ctx:                      config.ctx,
		VerifierEmail:            config.verifierEmail,
		VerifierDomain:           config.verifierDomain,
		ValidationTypeDefault:    config.validationTypeDefault,
		ConnectionTimeout:        config.connectionTimeout,
		ResponseTimeout:          config.responseTimeout,
		ConnectionAttempts:       config.connectionAttempts,
		WhitelistedDomains:       config.whitelistedDomains,
		BlacklistedDomains:       config.blacklistedDomains,
		BlacklistedMxIpAddresses: config.blacklistedMxIpAddresses,
		DNS:                      dns,
		ValidationTypeByDomain:   config.validationTypeByDomain,
		WhitelistValidation:      config.whitelistValidation,
		NotRfcMxLookupFlow:       config.notRfcMxLookupFlow,
		SMTPFailFast:             config.smtpFailFast,
		SMTPSafeCheck:            config.smtpSafeCheck,
		EmailPattern:             regexEmail,
		SMTPErrorBodyPattern:     regexSMTPErrorBody,
	}
	return &newConfiguration, err
}
