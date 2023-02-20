# ![Truemail - configurable framework agnostic plain Go email validator](https://truemail-rb.org/assets/images/truemail_logo.png)

[![Go Report Card](https://goreportcard.com/badge/github.com/truemail-rb/truemail-go)](https://goreportcard.com/report/github.com/truemail-rb/truemail-go)
[![Codecov](https://codecov.io/gh/truemail-rb/truemail-go/branch/master/graph/badge.svg)](https://codecov.io/gh/truemail-rb/truemail-go)
[![CircleCI](https://circleci.com/gh/truemail-rb/truemail-go/tree/master.svg?style=svg)](https://circleci.com/gh/truemail-rb/truemail-go/tree/master)
[![GitHub release (latest by date)](https://img.shields.io/github/v/release/truemail-rb/truemail-go)](https://github.com/truemail-rb/truemail-go/releases)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/truemail-rb/truemail-go)](https://pkg.go.dev/github.com/truemail-rb/truemail-go)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go)
[![GitHub](https://img.shields.io/github/license/truemail-rb/truemail-go)](LICENSE.txt)
[![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-v1.4%20adopted-ff69b4.svg)](CODE_OF_CONDUCT.md)

Configurable Golang 📨 email validator. Verify email via Regex, DNS, SMTP and even more. Be sure that email address valid and exists. It's Golang port of [`truemail`](https://truemail-rb.org/truemail-gem) Ruby gem. Currently ported all validation features only.

> Actual and maintainable documentation :books: for developers is living [here](https://truemail-rb.org/truemail-go).

## Table of Contents

- [Synopsis](#synopsis)
- [Features](#features)
- [Requirements](#requirements)
- [Installation](#installation)
- [Usage](#usage)
  - [Configuration features](#configuration-features)
    - [Creating configuration](#creating-configuration)
    - [Using configuration](#using-configuration)
  - [Validation features](#validation-features)
    - [Whitelist/Blacklist check](#whitelistblacklist-check)
      - [Whitelist case](#whitelist-case)
      - [Whitelist validation case](#whitelist-validation-case)
      - [Blacklist case](#blacklist-case)
      - [Duplication case](#duplication-case)
    - [Regex validation](#regex-validation)
      - [With default regex pattern](#with-default-regex-pattern)
      - [With custom regex pattern](#with-custom-regex-pattern)
    - [DNS (MX) validation](#mx-validation)
      - [RFC MX lookup flow](#rfc-mx-lookup-flow)
      - [Not RFC MX lookup flow](#not-rfc-mx-lookup-flow)
    - [MX blacklist validation](#mx-blacklist-validation)
    - [SMTP validation](#smtp-validation)
      - [SMTP fail fast enabled](#smtp-fail-fast-enabled)
      - [SMTP safe check disabled](#smtp-safe-check-disabled)
      - [SMTP safe check enabled](#smtp-safe-check-enabled)
- [Truemail helpers](#truemail-helpers)
- [Truemail family](#truemail-family)
- [Contributing](#contributing)
- [License](#license)
- [Code of Conduct](#code-of-conduct)
- [Credits](#credits)
- [Versioning](#versioning)
- [Changelog](CHANGELOG.md)

## Synopsis

Email validation is a tricky thing. There are a number of different ways to validate an email address and all mechanisms must conform with the best practices and provide proper validation. The `truemail-go` package helps you validate emails via regex pattern, presence of DNS records, and real existence of email account on a current email server.

**Syntax Checking**: Checks the email addresses via regex pattern.

**Mail Server Existence Check**: Checks the availability of the email address domain using DNS records.

**Mail Existence Check**: Checks if the email address really exists and can receive email via SMTP connections and email-sending emulation techniques.

## Features

- Configurable validator, validate only what you need
- Supporting of internationalized emails ([EAI](https://en.wikipedia.org/wiki/Email_address#Internationalization))
- Whitelist/blacklist validation layers
- Ability to configure different MX/SMTP validation flows
- Ability to configure [DEA](https://en.wikipedia.org/wiki/Disposable_email_address) validation flow
- Simple SMTP debugger

## Requirements

Golang 1.19+

## Installation

Install `truemail-go`:

```bash
go get github.com/truemail-rb/truemail-go
go install -i github.com/truemail-rb/truemail-go
```

Import `truemail-go` dependency into your code:

```go
package main

import "github.com/truemail-rb/truemail-go"
```

## Usage

### Configuration features

You can use global package configuration or custom independent configuration. Available configuration options:

- verifier email
- verifier domain
- email pattern
- connection timeout
- response timeout
- connection attempts
- default validation type
- validation type for domains
- whitelisted domains
- whitelist validation
- blacklisted domains
- blacklisted mx ip-addresses
- custom DNS gateway
- RFC MX lookup flow
- SMTP port number
- SMTP error body pattern
- SMTP fail fast
- SMTP safe check

#### Creating configuration

To have an access for library features, you must create configuration struct first. Please use `truemail.NewConfiguration()` built-in constructor to create a valid configuration as in the example below:

```go
import "github.com/truemail-rb/truemail-go"

configuration := truemail.NewConfiguration(
  ConfigurationAttr{
    // Required parameter. Must be an existing email on behalf of which verification will be
    // performed
    VerifierEmail: "verifier@example.com",

    // Optional parameter. Must be an existing domain on behalf of which verification will be
    // performed. By default verifier domain based on verifier email
    VerifierDomain: "somedomain.com",

    // Optional parameter. You can override default regex pattern
    EmailPattern: `\A.+@(.+)\z`,

    // Optional parameter. You can override default regex pattern
    SmtpErrorBodyPattern: `.*(user|account).*`,

    // Optional parameter. Connection timeout in seconds.
    // It is equal to 2 by default.
    ConnectionTimeout: 1,

    // Optional parameter. A SMTP server response timeout in seconds.
    // It is equal to 2 by default.
    ResponseTimeout: 1,

    // Optional parameter. Total of connection attempts. It is equal to 2 by default.
    // This parameter uses in mx lookup timeout error and smtp request (for cases when
    // there is one mx server).
    ConnectionAttempts: 1,

    // Optional parameter. You can predefine default validation type for
    // truemail.Validate("email@email.com", configuration) call without type-parameter
    // Available validation types: "regex", "mx", "mx_blacklist", "smtp"
    ValidationTypeDefault: "mx",

    // Optional parameter. You can predefine which type of validation will be used for domains.
    // Also you can skip validation by domain.
    // Available validation types: "regex", "mx", "smtp"
    // This configuration will be used over current or default validation type parameter
    // All of validations for "somedomain.com" will be processed with regex validation only.
    // And all of validations for "otherdomain.com" will be processed with mx validation only.
    // It is equal to empty map of strings by default.
    ValidationTypeByDomain: map[string]string{"somedomain.com": "regex", "otherdomain.com": "mx"},

    // Optional parameter. Validation of email which contains whitelisted domain always will
    // return true. Other validations will not processed even if it was defined in
    // validationTypeByDomain. It is equal to empty slice of strings by default.
    WhitelistedDomains: []string{"somedomain1.com", "somedomain2.com"},

    // Optional parameter. With this option Truemail will validate email which contains whitelisted
    // domain only, i.e. if domain whitelisted, validation will passed to Regex, MX or SMTP
    // validators. Validation of email which not contains whitelisted domain always will return
    // false. It is equal false by default.
    WhitelistValidation: true,

    // Optional parameter. Validation of email which contains blacklisted domain always will
    // return false. Other validations will not processed even if it was defined in
    // validationTypeByDomain. It is equal to empty slice of strings by default.
    BlacklistedDomains: []string{"somedomain3.com", "somedomain4.com"},

    // Optional parameter. With this option Truemail will filter out unwanted mx servers via
    // predefined list of ip addresses. It can be used as a part of DEA (disposable email
    // address) validations. It is equal to empty slice of strings by default.
    BlacklistedMxIpAddresses: []string{"1.1.1.1", "2.2.2.2"},

    // Optional parameter. This option will provide to use custom DNS gateway when Truemail
    // interacts with DNS. Valid port number is in the range 1-65535. If you won't specify
    // nameserver port Truemail will use default DNS TCP/UDP port 53. It means that you can
    // use ip4 addres as DNS gateway, for example "10.0.0.1". By default Truemail uses
    // DNS gateway from system settings and this option is equal to empty string.
    Dns: "10.0.0.1:5300",

    // Optional parameter. This option will provide to use not RFC MX lookup flow.
    // It means that MX and Null MX records will be cheked on the DNS validation layer only.
    // By default this option is disabled and equal to false.
    NotRfcMxLookupFlow: true,

    // Optional parameter. SMTP port number. It is equal to 25 by default.
    // This parameter uses for SMTP session in SMTP validation layer.
    SmtpPort: 2525,

    // Optional parameter. This option will provide to use smtp fail fast behavior. When
    // smtpFailFast = true it means that Truemail ends smtp validation session after first
    // attempt on the first mx server in any fail cases (network connection/timeout error,
    // smtp validation error). This feature helps to reduce total time of SMTP validation
    // session up to 1 second. By default this option is disabled and equal to false.
    SmtpFailFast: true,

    // Optional parameter. This option will be parse bodies of SMTP errors. It will be helpful
    // if SMTP server does not return an exact answer that the email does not exist
    // By default this option is disabled, available for SMTP validation only.
    SmtpSafeCheck: true,
  },
)
```

#### Using configuration

```go
import "github.com/truemail-rb/truemail-go"

configuration := truemail.NewConfiguration(truemail.ConfigurationAttr{VerifierEmail: "verifier@example.com"})

truemail.Validate("some@email.com", configuration)
truemail.IsValid("some@email.com", configuration, "regex")
```

### Validation features

#### Whitelist/Blacklist check

Whitelist/Blacklist check is zero validation level. You can define white and black list domains. It means that validation of email which contains whitelisted domain always will return `true`, and for blacklisted domain will return `false`.

Please note, other validations will not processed even if it was defined in `ValidationTypeByDomain`.

**Sequence of domain list check:**

1. Whitelist check
2. Whitelist validation check
3. Blacklist check

Example of usage:

```go
import "github.com/truemail-rb/truemail-go"

configuration := truemail.NewConfiguration(
  truemail.ConfigurationAttr{
    VerifierEmail: "verifier@example.com",
    WhitelistedDomains: []string{"white-domain.com", "somedomain.com"},
    BlacklistedDomains: []string{"black-domain.com", "somedomain.com"},
    ValidationTypeByDomain: map[string]string{"somedomain.com": "mx"},
  },
)
```

##### Whitelist case

When email in whitelist, validation type will be redefined. Validation result returns ```true```

```go
import "github.com/truemail-rb/truemail-go"

configuration := truemail.NewConfiguration(
  truemail.ConfigurationAttr{
    VerifierEmail: "verifier@example.com",
    WhitelistedDomains: []string{"white-domain.com", "somedomain.com"},
    BlacklistedDomains: []string{"black-domain.com", "somedomain.com"},
    ValidationTypeByDomain: map[string]string{"somedomain.com": "mx"},
  },
)

truemail.Validate("email@white-domain.com", configuration) // returns pointer to ValidatorResult with validation details and error
truemail.IsValid("email@white-domain.com", configuration) // returns true
```

##### Whitelist validation case

When email domain in whitelist and `WhitelistValidation` is sets equal to `true` validation type will be passed to other validators. Validation of email which not contains whitelisted domain always will return `false`.

###### Email has whitelisted domain

```go
import "github.com/truemail-rb/truemail-go"

configuration := truemail.NewConfiguration(
  truemail.ConfigurationAttr{
    VerifierEmail: "verifier@example.com",
    WhitelistedDomains: []string{"white-domain.com"},
    WhitelistValidation: true,
  },
)

truemail.Validate("email@white-domain.com", configuration, "regex") // returns pointer to ValidatorResult with validation details and error
truemail.IsValid("email@white-domain.com", configuration, "regex") // returns true
```

###### Email hasn't whitelisted domain

```go
import "github.com/truemail-rb/truemail-go"

configuration := truemail.NewConfiguration(
  truemail.ConfigurationAttr{
    VerifierEmail: "verifier@example.com",
    WhitelistedDomains: []string{"white-domain.com"},
    WhitelistValidation: true,
  },
)

truemail.Validate("email@domain.com", configuration, "regex") // returns pointer to ValidatorResult with validation details and error
truemail.IsValid("email@domain.com", configuration, "regex") // returns false
```

##### Blacklist case

When email in blacklist, validation type will be redefined too. Validation result returns `false`.

```go
import "github.com/truemail-rb/truemail-go"

configuration := truemail.NewConfiguration(
  truemail.ConfigurationAttr{
    VerifierEmail: "verifier@example.com",
    WhitelistedDomains: []string{"white-domain.com", "somedomain.com"},
    BlacklistedDomains: []string{"black-domain.com", "somedomain.com"},
    ValidationTypeByDomain: map[string]string{"somedomain.com": "mx"},
  },
)

truemail.Validate("email@black-domain.com", configuration) // returns pointer to ValidatorResult with validation details and error
truemail.IsValid("email@domain.com", configuration) // returns false
```

##### Duplication case

Validation result for this email returns `true`, because it was found in whitelisted domains list first. Also `ValidatorResult.ValidationType` for this case will be redefined.

```go
import "github.com/truemail-rb/truemail-go"

configuration := truemail.NewConfiguration(
  truemail.ConfigurationAttr{
    VerifierEmail: "verifier@example.com",
    WhitelistedDomains: []string{"white-domain.com", "somedomain.com"},
    BlacklistedDomains: []string{"black-domain.com", "somedomain.com"},
    ValidationTypeByDomain: map[string]string{"somedomain.com": "mx"},
  },
)

truemail.Validate("email@somedomain.com", configuration) // returns pointer to ValidatorResult with validation details and error
truemail.IsValid("email@somedomain.com", configuration) // returns true
```

#### Regex validation

Validation with regex pattern is the first validation level. It uses whitelist/blacklist check before running itself.

```code
[Whitelist/Blacklist] -> [Regex validation]
```

By default this validation not performs strictly following [RFC 5322](https://www.ietf.org/rfc/rfc5322.txt) standard, so you can override `truemail` default regex pattern if you want.

Example of usage:

##### With default regex pattern

```go
import "github.com/truemail-rb/truemail-go"

configuration := truemail.NewConfiguration(
  truemail.ConfigurationAttr{
    VerifierEmail: "verifier@example.com",
  },
)

truemail.Validate("email@example.com", configuration, "regex") // returns pointer to ValidatorResult with validation details and error
truemail.IsValid("email@example.com", configuration, "regex") // returns true
```

##### With custom regex pattern

You should define your custom regex pattern in a gem configuration before.

```go
import "github.com/truemail-rb/truemail-go"

configuration := truemail.NewConfiguration(
  truemail.ConfigurationAttr{
    VerifierEmail: "verifier@example.com",
    EmailPattern: `\A(.+)@(.+)\z`,
  },
)

truemail.Validate("email@example.com", configuration, "regex") // returns pointer to ValidatorResult with validation details and error
truemail.IsValid("email@example.com", configuration, "regex") // returns true
truemail.IsValid("not_email", configuration, "regex") // returns false
```

#### MX validation

In fact it's DNS validation because it checks not MX records only. DNS validation is the second validation level, historically named as MX validation. It uses Regex validation before running itself. When regex validation has completed successfully then runs itself.

```code
[Whitelist/Blacklist] -> [Regex validation] -> [MX validation]
```

Please note, `truemail` MX validator [not performs](https://github.com/truemail-rb/truemail-go/issues/26) strict compliance of the [RFC 5321](https://tools.ietf.org/html/rfc5321#section-5) standard for best validation outcome.

##### RFC MX lookup flow

[Truemail MX lookup](https://slides.com/vladislavtrotsenko/truemail#/0/9) based on RFC 5321. It consists of 3 substeps: MX, CNAME and A record resolvers. The point of each resolver is attempt to extract the mail servers from email domain. If at least one server exists that validation is successful. Iteration is processing until resolver returns true.

Example of usage:

```go
import "github.com/truemail-rb/truemail-go"

configuration := truemail.NewConfiguration(
  truemail.ConfigurationAttr{
    VerifierEmail: "verifier@example.com",
  },
)

truemail.Validate("email@example.com", configuration, "mx") // returns pointer to ValidatorResult with validation details and error
truemail.IsValid("email@example.com", configuration, "mx") // returns bool
```

##### Not RFC MX lookup flow

Also Truemail has possibility to use not RFC MX lookup flow. It means that will be used only one MX resolver on the DNS validation layer. By default this option is disabled.

Example of usage:

```go
import "github.com/truemail-rb/truemail-go"

configuration := truemail.NewConfiguration(
  truemail.ConfigurationAttr{
    VerifierEmail: "verifier@example.com",
    NotRfcMxLookupFlow: true,
  },
)

truemail.Validate("email@example.com", configuration, "mx") // returns pointer to ValidatorResult with validation details and error
truemail.IsValid("email@example.com", configuration, "mx") // returns bool
```

#### MX blacklist validation

MX blacklist validation is the third validation level. This layer provides checking extracted mail server(s) IP address from MX validation with predefined blacklisted IP addresses list. It can be used as a part of DEA ([disposable email address](https://en.wikipedia.org/wiki/Disposable_email_address)) validations.

```code
[Whitelist/Blacklist] -> [Regex validation] -> [MX validation] -> [MX blacklist validation]
```

Example of usage:

```go
import "github.com/truemail-rb/truemail-go"

configuration := truemail.NewConfiguration(
  truemail.ConfigurationAttr{
    VerifierEmail: "verifier@example.com",
    BlacklistedMxIpAddresses: []string{"127.0.1.2", "127.0.1.3"},
  },
)

truemail.Validate("email@example.com", configuration, "mx_blacklist") // returns pointer to ValidatorResult with validation details and error
truemail.IsValid("email@example.com", configuration, "mx_blacklist") // returns bool
```

#### SMTP validation

SMTP validation is a final, fourth validation level. This type of validation tries to check real existence of email account on a current email server. This validation runs a chain of previous validations and if they're complete successfully then runs itself.

```code
[Whitelist/Blacklist] -> [Regex validation] -> [MX validation] -> [MX blacklist validation] -> [SMTP validation]
```

If total count of MX servers is equal to one, `truemail` SMTP validator will use value from `ConnectionAttempts` as connection attempts. By default it's equal `2`.

By default, you don't need pass with-parameter to use it. Example of usage is specified below:

##### SMTP fail fast enabled

Truemail can use fail fast behavior for SMTP validation layer. When `SmtpFailFast = true` it means that `truemail` ends smtp validation session after first attempt on the first mx server in any fail cases (network connection/timeout error, smtp validation error). This feature helps to reduce total time of SMTP validation session up to 1 second.

```go
import "github.com/truemail-rb/truemail-go"

configuration := truemail.NewConfiguration(
  truemail.ConfigurationAttr{
    VerifierEmail: "verifier@example.com",
    SmtpFailFast: true,
  },
)

truemail.Validate("email@example.com", configuration) // returns pointer to ValidatorResult with validation details and error
truemail.IsValid("email@example.com", configuration) // returns bool
```

##### SMTP safe check disabled

When this feature disabled, it means that SMTP validation will be failed when it consists at least one smtp error.

```go
import "github.com/truemail-rb/truemail-go"

configuration := truemail.NewConfiguration(
  truemail.ConfigurationAttr{
    VerifierEmail: "verifier@example.com",
  },
)

truemail.Validate("email@example.com", configuration) // returns pointer to ValidatorResult with validation details and error
truemail.IsValid("email@example.com", configuration) // returns bool
```

##### SMTP safe check enabled

When this feature enabled, it means that SMTP validation will be successful for all cases until `truemail` SMTP validator receive `RCPT TO` error that matches to `SmtpErrorBodyPattern`, specified in `configuration`.

```go
import "github.com/truemail-rb/truemail-go"

configuration := truemail.NewConfiguration(
  truemail.ConfigurationAttr{
    VerifierEmail: "verifier@example.com",
    SmtpSafeCheck: true,
  },
)

truemail.Validate("email@example.com", configuration) // returns pointer to ValidatorResult with validation details and error
truemail.IsValid("email@example.com", configuration) // returns bool
```

### Truemail helpers

#### .IsValid()

You can use the `.IsValid()` helper for quick validation of email address. It returns a boolean:

```go
truemail.IsValid("email@example.com", configuration)
```

## Truemail family

All Truemail solutions: <https://truemail-rb.org>

| Name | Type | Description |
| --- | --- | --- |
| [truemail](https://github.com/truemail-rb/truemail) | ruby gem | Configurable framework agnostic plain Ruby email validator, main core |
| [truemail server](https://github.com/truemail-rb/truemail-go-rack) | ruby app | Lightweight rack based web API wrapper for Truemail gem |
| [truemail-rack-docker](https://github.com/truemail-rb/truemail-go-rack-docker-image) | docker image | Lightweight rack based web API [dockerized image](https://hub.docker.com/r/truemail/truemail-rack) :whale: of Truemail server |
| [truemail-ruby-client](https://github.com/truemail-rb/truemail-go-ruby-client) | ruby gem | Web API Ruby client for Truemail Server |
| [truemail-crystal-client](https://github.com/truemail-rb/truemail-go-crystal-client) | crystal shard | Web API Crystal client for Truemail Server |
| [truemail-java-client](https://github.com/truemail-rb/truemail-go-java-client) | java lib | Web API Java client for Truemail Server |
| [truemail-rspec](https://github.com/truemail-rb/truemail-go-rspec) | ruby gem | Truemail configuration, auditor and validator RSpec helpers |

## Contributing

Bug reports and pull requests are welcome on GitHub at <https://github.com/truemail-rb/truemail-go>. This project is intended to be a safe, welcoming space for collaboration, and contributors are expected to adhere to the [Contributor Covenant](http://contributor-covenant.org) code of conduct. Please check the [open tickets](https://github.com/truemail-rb/truemail-go/issues). Be sure to follow Contributor Code of Conduct below and our [Contributing Guidelines](CONTRIBUTING.md).

## License

The package is available as open source under the terms of the [MIT License](https://opensource.org/licenses/MIT).

## Code of Conduct

Everyone interacting in the Truemail project’s codebases, issue trackers, chat rooms and mailing lists is expected to follow the [code of conduct](CODE_OF_CONDUCT.md).

## Credits

- [The Contributors](https://github.com/truemail-rb/truemail-go/graphs/contributors) for code and awesome suggestions
- [The Stargazers](https://github.com/truemail-rb/truemail-go/stargazers) for showing their support

## Versioning

Truemail uses [Semantic Versioning 2.0.0](https://semver.org)
