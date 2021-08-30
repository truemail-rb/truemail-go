package truemail

const (
	// network configuration options

	DefaultConnectionTimeout  = 2
	DefaultResponseTimeout    = 2
	DefaultConnectionAttempts = 2
	DefaultDnsPort            = 53
	DefaultSmtpPort           = 25
	TcpTransportLayer         = "tcp"

	// validation types

	ValidationTypeDomainListMatch = "domain_list_match"
	ValidationTypeRegex           = "regex"
	ValidationTypeMx              = "mx"
	ValidationTypeMxBlacklist     = "mx_blacklist"
	ValidationTypeSmtp            = "smtp"
	ValidationTypeDefault         = ValidationTypeSmtp

	// regex patterns

	DomainCharsSize              = `\A.{4,255}\z`
	EmailCharsSize               = `\A.{6,255}\z`
	RegexDomainPattern           = `(?i)[\p{L}0-9]+([\-.]{1}[\p{L}0-9]+)*\.\p{L}{2,63}`
	RegexEmailPattern            = `(\A([\p{L}0-9]+[\W\w]*)@(` + RegexDomainPattern + `)\z)`
	RegexDomainFromEmail         = `\A.+@(.+)\z`
	RegexSMTPErrorBodyPattern    = `(?i).*550{1}.*(user|account|customer|mailbox).*`
	RegexPortNumber              = `(6553[0-5]|655[0-2]\d|65[0-4](\d){2}|6[0-4](\d){3}|[1-5](\d){4}|[1-9](\d){0,3})`
	RegexIpAddress               = `((\d|[1-9]\d|1\d{2}|2[0-4]\d|25[0-5])\.){3}(\d|[1-9]\d|1\d{2}|2[0-4]\d|25[0-5])`
	RegexIpAddressPattern        = `\A` + RegexIpAddress + `\z`
	RegexDNSServerAddressPattern = `\A` + RegexIpAddress + `(:` + RegexPortNumber + `)?\z`

	// shortcuts

	EmptyString = ""

	// validationDomainListMatch

	DomainListMatchWhitelist    = "whitelist"
	DomainListMatchBlacklist    = "blacklist"
	DomainListMatchErrorContext = "blacklisted email"

	// validationRegex

	RegexErrorContext = "email does not match the regular expression"

	// validationMxBlacklist

	MxBlacklistErrorContext = "blacklisted mx server ip address"

	// validationMx

	MxErrorContext = "target host(s) not found"

	// validatorSmtp

	SmtpErrorContext = "smtp error"
)
