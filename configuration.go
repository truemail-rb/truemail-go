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

// NewConfiguration builder. Returns valid newConfiguration structure
func NewConfiguration(config ConfigurationAttr) (*configuration, error) {
	config.assignDefaultValues()
	err := config.validate()

	if err != nil {
		return nil, err
	}

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
		DNS:                      config.dns,
		ValidationTypeByDomain:   config.validationTypeByDomain,
		WhitelistValidation:      config.whitelistValidation,
		NotRfcMxLookupFlow:       config.notRfcMxLookupFlow,
		SMTPFailFast:             config.smtpFailFast,
		SMTPSafeCheck:            config.smtpSafeCheck,
		EmailPattern:             config.regexEmail,
		SMTPErrorBodyPattern:     config.regexSMTPErrorBody,
	}
	return &newConfiguration, err
}
