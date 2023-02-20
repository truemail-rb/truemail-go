package truemail

import (
	"context"
	"regexp"
)

// Configuration structure
type Configuration struct {
	ctx                                                                  context.Context
	VerifierEmail, VerifierDomain, ValidationTypeDefault, Dns            string
	ConnectionTimeout, ResponseTimeout, ConnectionAttempts, SmtpPort     int
	WhitelistedDomains, BlacklistedDomains, BlacklistedMxIpAddresses     []string
	ValidationTypeByDomain                                               map[string]string
	WhitelistValidation, NotRfcMxLookupFlow, SmtpFailFast, SmtpSafeCheck bool
	EmailPattern, SmtpErrorBodyPattern                                   *regexp.Regexp
}

// NewConfiguration returns new valid newConfiguration structure
func NewConfiguration(config ConfigurationAttr) (*Configuration, error) {
	config.assignDefaultValues()
	err := config.validate()

	if err != nil {
		return nil, err
	}

	newConfiguration := Configuration{
		ctx:                      config.ctx,
		VerifierEmail:            config.VerifierEmail,
		VerifierDomain:           config.VerifierDomain,
		ValidationTypeDefault:    config.ValidationTypeDefault,
		ConnectionTimeout:        config.ConnectionTimeout,
		ResponseTimeout:          config.ResponseTimeout,
		ConnectionAttempts:       config.ConnectionAttempts,
		WhitelistedDomains:       config.WhitelistedDomains,
		BlacklistedDomains:       config.BlacklistedDomains,
		BlacklistedMxIpAddresses: config.BlacklistedMxIpAddresses,
		Dns:                      config.Dns,
		ValidationTypeByDomain:   config.ValidationTypeByDomain,
		WhitelistValidation:      config.WhitelistValidation,
		NotRfcMxLookupFlow:       config.NotRfcMxLookupFlow,
		SmtpPort:                 config.SmtpPort,
		SmtpFailFast:             config.SmtpFailFast,
		SmtpSafeCheck:            config.SmtpSafeCheck,
		EmailPattern:             config.RegexEmail,
		SmtpErrorBodyPattern:     config.RegexSmtpErrorBody,
	}
	return &newConfiguration, err
}
