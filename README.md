# ![Truemail - configurable framework agnostic plain Go email validator](https://truemail-rb.org/assets/images/truemail_logo.png)

Configurable framework agnostic plain Go email validator. Verify email via Regex, DNS, SMTP and even more. Be sure that email address valid and exists.

> Actual and maintainable documentation :books: for developers is living [here](https://truemail-rb.org/truemail-go).

## Table of Contents

- [Synopsis](#synopsis)
- [Features](#features)
- [Requirements](#requirements)
- [Installation](#installation)
- [Usage](#usage)

## Synopsis

Email validation is a tricky thing. There are a number of different ways to validate an email address and all mechanisms must conform with the best practices and provide proper validation. The Truemail library helps you validate emails via regex pattern, presence of DNS records, and real existence of email account on a current email server.

**Syntax Checking**: Checks the email addresses via regex pattern.

**Mail Server Existence Check**: Checks the availability of the email address domain using DNS records.

**Mail Existence Check**: Checks if the email address really exists and can receive email via SMTP connections and email-sending emulation techniques.

Also Truemail library allows performing an audit of the host in which runs.

## Features

- Configurable validator, validate only what you need
- Supporting of internationalized emails ([EAI](https://en.wikipedia.org/wiki/Email_address#Internationalization))
- Whitelist/blacklist validation layers
- Ability to configure different MX/SMTP validation flows
- Ability to configure [DEA](https://en.wikipedia.org/wiki/Disposable_email_address) validation flow
- Simple SMTP debugger
- Event logger
- Host auditor tools (helps to detect common host problems interfering to proper email verification)
- JSON serializers

## Requirements

Golang 1.15+

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

You can use global gem configuration or custom independent configuration. Available configuration options:

- verifier email
- verifier domain
- email pattern
- SMTP error body pattern
- connection timeout
- response timeout
- connection attempts
- default validation type
- validation type for domains
- whitelisted domains
- whitelist validation
- blacklisted domains
- blacklisted mx ip-addresses
- custom DNS gateway(s)
- RFC MX lookup flow
- SMTP fail fast
- SMTP safe check
- JSON serializer

#### Creating configuration

To have an access for library features, you must create configuration struct first. Please use `truemail.NewConfiguration()` built-in constructor to create a valid configuration as in the example below:

```go
import "github.com/truemail-rb/truemail-go"

configuration := truemail.NewConfiguration(
  ConfigurationAttr{
    // Required parameter. Must be an existing email on behalf of which verification will be
    // performed
    verifierEmail: "verifier@example.com",

    // Optional parameter. Must be an existing domain on behalf of which verification will be
    // performed. By default verifier domain based on verifier email
    verifierDomain: "somedomain.com",

    // Optional parameter. You can override default regex pattern
    emailPattern: `\A.+@(.+)\z`,

    // Optional parameter. You can override default regex pattern
    smtpErrorBodyPattern: `.*(user|account).*`,

    // Optional parameter. Connection timeout in seconds.
    // It is equal to 2 by default.
    connectionTimeout: 1,

    // Optional parameter. A SMTP server response timeout in seconds.
    // It is equal to 2 by default.
    responseTimeout: 1,

    // Optional parameter. Total of connection attempts. It is equal to 2 by default.
    // This parameter uses in mx lookup timeout error and smtp request (for cases when
    // there is one mx server).
    connectionAttempts: 1,

    // Optional parameter. You can predefine default validation type for
    // truemail.Validate("email@email.com", configuration) call without type-parameter
    // Available validation types: "regex", "mx", "mx_blacklist", "smtp"
    validationTypeDefault: "mx",

    // Optional parameter. You can predefine which type of validation will be used for domains.
    // Also you can skip validation by domain.
    // Available validation types: "regex", "mx", "smtp"
    // This configuration will be used over current or default validation type parameter
    // All of validations for "somedomain.com" will be processed with regex validation only.
    // And all of validations for "otherdomain.com" will be processed with mx validation only.
    // It is equal to empty map of strings by default.
    validationTypeByDomain: map[string]string{"somedomain.com": "regex", "otherdomain.com": "mx"},

    // Optional parameter. Validation of email which contains whitelisted domain always will
    // return true. Other validations will not processed even if it was defined in
    // validationTypeByDomain. It is equal to empty slice of strings by default.
    whitelistedDomains: []string{"somedomain1.com", "somedomain2.com"},

    // Optional parameter. With this option Truemail will validate email which contains whitelisted
    // domain only, i.e. if domain whitelisted, validation will passed to Regex, MX or SMTP
    // validators. Validation of email which not contains whitelisted domain always will return
    // false. It is equal false by default.
    whitelistValidation: true,

    // Optional parameter. Validation of email which contains blacklisted domain always will
    // return false. Other validations will not processed even if it was defined in
    // validationTypeByDomain. It is equal to empty slice of strings by default.
    blacklistedDomains: []string{"somedomain3.com", "somedomain4.com"},

    // Optional parameter. With this option Truemail will filter out unwanted mx servers via
    // predefined list of ip addresses. It can be used as a part of DEA (disposable email
    // address) validations. It is equal to empty slice of strings by default.
    blacklistedMxIpAddresses: []string{"1.1.1.1", "2.2.2.2"},

    // Optional parameter. This option will provide to use custom DNS gateway when Truemail
    // interacts with DNS. Valid port numbers are in the range 1-65535. If you won't specify
    // nameserver's ports Truemail will use default DNS TCP/UDP port 53. By default Truemail
    // uses DNS gateway from system settings and this option is equal to empty slice of
    // strings by default.
    dns: []string{"10.0.0.1", "10.0.0.2:54"},

    // Optional parameter. This option will provide to use not RFC MX lookup flow.
    // It means that MX and Null MX records will be cheked on the DNS validation layer only.
    // By default this option is disabled and equal to false.
    notRfcMxLookupFlow: true,

    // Optional parameter. This option will provide to use smtp fail fast behaviour. When
    // smtpFailFast = true it means that Truemail ends smtp validation session after first
    // attempt on the first mx server in any fail cases (network connection/timeout error,
    // smtp validation error). This feature helps to reduce total time of SMTP validation
    // session up to 1 second. By default this option is disabled and equal to false.
    smtpFailFast: true,

    // Optional parameter. This option will be parse bodies of SMTP errors. It will be helpful
    // if SMTP server does not return an exact answer that the email does not exist
    // By default this option is disabled, available for SMTP validation only.
    smtpSafeCheck: true,
  },
)
```

#### Using configuration

```go
import "github.com/truemail-rb/truemail-go"

configuration := truemail.NewConfiguration(ConfigurationAttr{verifierEmail: "verifier@example.com"})

truemail.Validate("some@email.com", configuration)
truemail.IsValid("some@email.com", configuration, "regex")
```
