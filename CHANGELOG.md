# Changelog

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.5] - 2023-02-20

### Updated

- Updated `x/crypto` dependency

### Fixed

- Fixed `x/crypto` issue. Vulnerable to panic via SSH server, [CVE-2021-43565](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2021-43565)
- Fixed `x/crypto` issue. Use of a Broken or Risky Cryptographic Algorithm in golang.org/x/crypto/ssh, [CVE-2022-27191](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2022-27191)
- Fixed `x/crypto` issue. Panic in malformed certificate, [CVE-2020-7919](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2020-7919)
- Fixed `x/crypto` issue. Improper Verification of Cryptographic Signature, [CVE-2020-9283](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2020-9283)

## [1.0.4] - 2023-02-19

### Added

- Added tag script for new release tagging

### Updated

- Updated `go` to `1.20.1`
- Updated `x/net` dependency
- Updated releasing script (auto deploy to `GitHub`/`Go Packages`)
- Updated `CircleCI` config
- Updated project license

### Fixed

- Updated `x/net` dependency. A maliciously crafted HTTP/2 stream could cause excessive CPU consumption in the `HPACK` decoder, sufficient to cause a denial of service from a small number of small requests, [CVE-2022-41723](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2022-41723)

## [1.0.3] - 2022-12-26

### Added

- Added [`cspell`](https://cspell.org) linter
- Added [`markdownlint`](https://github.com/DavidAnson/markdownlint) linter
- Added [`shellcheck`](https://www.shellcheck.net) linter
- Added [`yamllint`](https://yamllint.readthedocs.io) linter
- Added [`lefthook`](https://github.com/evilmartians/lefthook) linters aggregator

### Fixed

- Fixed typos in project's codebase
- Fixed new project's linter issues

### Updated

- Updated `CircleCI` config

## [1.0.2] - 2022-11-20

- Updated dependencies

## [1.0.1] - 2022-10-27

- Updated dependencies

## [1.0.0] - 2022-10-13

- Updated minimal go version to 1.19
- Updated dependencies
- Updated circleci config
- Updated package documentation

## [0.1.4] - 2022-05-31

### Fixed

- Updated `yaml.v3` indirect dependency. An issue in the Unmarshal function in Go-Yaml v3 causes the program to crash when attempting to deserialize invalid input, [CVE-2022-28948](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2022-28948)

## [0.1.3] - 2022-03-02

### Added

- Added changelog

### Updated

- Updated `ConfigurationAttr#validate`, tests
- Updated `golangci` config
- Updated package documentation

## [0.1.2] - 2022-02-28

### Fixed

- Fixed linters issues

## [0.1.1] - 2022-02-28

### Added

- Added codecov

## [0.1.0] - 2022-02-28

### Added

- First release of `truemail-go`.
