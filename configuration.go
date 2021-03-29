package truemail

import (
	"context"
	"regexp"
)

// Configuration structure
type Configuration struct {
	ctx                                                                  context.Context
	VerifierEmail, VerifierDomain, ValidationTypeDefault                 string
	ConnectionTimeout, ResponseTimeout, ConnectionAttempts               int
	WhitelistedDomains, BlacklistedDomains, DNS                          []string
	ValidationTypeByDomain                                               map[string]string
	WhitelistValidation, NotRfcMxLookupFlow, SMTPFailFast, SMTPSafeCheck bool
	EmailPattern, SMTPErrorBodyPattern                                   *regexp.Regexp
}

// ConfigurationAttr kwargs for configuration builder
type ConfigurationAttr struct {
	ctx                                                                                      context.Context
	verifierEmail, verifierDomain, ValidationTypeDefault, emailPattern, smtpErrorBodyPattern string
	connectionTimeout, responseTimeout, connectionAttempts                                   int
	whitelistedDomains, blacklistedDomains, dns                                              []string
	validationTypeByDomain                                                                   map[string]string
	whitelistValidation, notRfcMxLookupFlow, smtpFailFast, smtpSafeCheck                     bool
}

// Configuration builder constants, regex patterns
const (
	DefaultConnectionTimeout     = 2
	DefaultResponseTimeout       = 2
	DefaultConnectionAttempts    = 2
	ValidationTypeRegex          = "regex"
	ValidationTypeMx             = "mx"
	ValidationTypeSMTP           = "smtp"
	ValidationTypeDefault        = ValidationTypeSMTP
	DomainCharsSize              = `\A.{4,255}\z`
	EmailCharsSize               = `\A.{6,255}\z`
	RegexDomainPattern           = `(?i)[\p{L}0-9]+([\-.]{1}[\p{L}0-9]+)*\.\p{L}{2,63}`
	RegexEmailPattern            = `(\A([\p{L}0-9]+[\w|\-.+]*)@(` + RegexDomainPattern + `)\z)`
	RegexDomainFromEmail         = `\A.+@(.+)\z`
	RegexSMTPErrorBodyPattern    = `(?i).*550{1}.*(user|account|customer|mailbox).*`
	RegexPortNumber              = `(6553[0-5]|655[0-2][0-9]\d|65[0-4](\d){2}|6[0-4](\d){3}|[1-5](\d){4}|[1-9](\d){0,3})`
	RegexDNSServerAddressPattern = `\A((\d|[1-9]\d|1\d{2}|2[0-4]\d|25[0-5])\.){3}(\d|[1-9]\d|1\d{2}|2[0-4]\d|25[0-5])(:` + RegexPortNumber + `)?\z`
)

// NewConfiguration builder
func NewConfiguration(config ConfigurationAttr) (*Configuration, error) {
	// assign fileds default values
	if config.ValidationTypeDefault == "" {
		config.ValidationTypeDefault = ValidationTypeDefault
	}
	if config.emailPattern == "" {
		config.emailPattern = RegexEmailPattern
	}
	if config.smtpErrorBodyPattern == "" {
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

	err = validateValidationTypeDefaultContext(config.ValidationTypeDefault)
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

	err = validateDNSServersContext(config.dns)
	if err != nil {
		return nil, err
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

	// create new Configuration
	newConfiguration := Configuration{
		ctx:                    config.ctx,
		VerifierEmail:          config.verifierEmail,
		VerifierDomain:         config.verifierDomain,
		ValidationTypeDefault:  config.ValidationTypeDefault,
		ConnectionTimeout:      config.connectionTimeout,
		ResponseTimeout:        config.responseTimeout,
		ConnectionAttempts:     config.connectionAttempts,
		WhitelistedDomains:     config.whitelistedDomains,
		BlacklistedDomains:     config.blacklistedDomains,
		DNS:                    config.dns,
		ValidationTypeByDomain: config.validationTypeByDomain,
		WhitelistValidation:    config.whitelistValidation,
		NotRfcMxLookupFlow:     config.notRfcMxLookupFlow,
		SMTPFailFast:           config.smtpFailFast,
		SMTPSafeCheck:          config.smtpSafeCheck,
		EmailPattern:           regexEmail,
		SMTPErrorBodyPattern:   regexSMTPErrorBody,
	}
	return &newConfiguration, err
}
