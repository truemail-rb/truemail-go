# Changelog

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.4] - 2022-05-31

### Fixed

- Updated `yaml.v3` indirect dependency. An issue in the Unmarshal function in Go-Yaml v3 causes the program to crash when attempting to deserialize invalid input, [CVE-2022-28948](https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2022-28948)

## [0.1.3] - 2022-03-02

### Added

- Added changelog

### Updated

- Updated `ConfigurationAttr#validate`, tests
- Updated golangci config
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
