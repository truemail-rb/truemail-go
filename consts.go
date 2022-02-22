package truemail

const (
	// network configuration options

	defaultConnectionTimeout  = 2
	defaultResponseTimeout    = 2
	defaultConnectionAttempts = 2
	defaultDnsPort            = 53
	defaultSmtpPort           = 25
	tcpTransportLayer         = "tcp"

	// validation types

	validationTypeDomainListMatch = "domain_list_match"
	validationTypeRegex           = "regex"
	validationTypeMx              = "mx"
	validationTypeMxBlacklist     = "mx_blacklist"
	validationTypeSmtp            = "smtp"
	validationTypeDefault         = validationTypeSmtp

	// regex patterns

	domainCharsSize              = `\A.{4,255}\z`
	emailCharsSize               = `\A.{6,255}\z`
	regexDomainPattern           = `(?i)[\p{L}0-9]+([\-.]{1}[\p{L}0-9]+)*\.\p{L}{2,63}`
	regexEmailPattern            = `(\A([\p{L}0-9]+[\W\w]*)@(` + regexDomainPattern + `)\z)`
	regexDomainFromEmail         = `\A.+@(.+)\z`
	regexSMTPErrorBodyPattern    = `(?i).*550{1}.*(user|account|customer|mailbox).*`
	regexPortNumber              = `(6553[0-5]|655[0-2]\d|65[0-4](\d){2}|6[0-4](\d){3}|[1-5](\d){4}|[1-9](\d){0,3})`
	regexIpAddress               = `((\d|[1-9]\d|1\d{2}|2[0-4]\d|25[0-5])\.){3}(\d|[1-9]\d|1\d{2}|2[0-4]\d|25[0-5])`
	regexIpAddressPattern        = `\A` + regexIpAddress + `\z`
	regexDNSServerAddressPattern = `\A` + regexIpAddress + `(:` + regexPortNumber + `)?\z`

	// shortcuts

	emptyString = ""

	// validationDomainListMatch

	domainListMatchWhitelist    = "whitelist"
	domainListMatchBlacklist    = "blacklist"
	domainListMatchErrorContext = "blacklisted email"

	// validationRegex

	regexErrorContext = "email does not match the regular expression"

	// validationMxBlacklist

	mxBlacklistErrorContext = "blacklisted mx server ip address"

	// validationMx

	mxErrorContext = "target host(s) not found"

	// validatorSmtp

	smtpErrorContext = "smtp error"
)
