package truemail

import (
	"fmt"
	"regexp"
)

// package helpers functions

// Returns true if the given string is present in slice,
// otherwise returns false.
func isIncluded(slice []string, target string) bool {
	for _, item := range slice {
		if item == target {
			return true
		}
	}

	return false
}

// Returns true if string from target slice found in base slice,
// otherwise returns false.
func isIntersected(baseSlice []string, targetSlice []string) bool {
	for _, item := range targetSlice {
		if isIncluded(baseSlice, item) {
			return true
		}
	}

	return false
}

// Regex builder
func newRegex(regexPattern string) (*regexp.Regexp, error) {
	return regexp.Compile(regexPattern)
}

// Matches string to regex pattern
func matchRegex(strContext, regexPattern string) bool {
	regex, err := newRegex(regexPattern)
	if err != nil {
		return false
	}

	return regex.MatchString(strContext)
}

// Returns slice of available validation types
func availableValidationTypes() []string {
	return []string{validationTypeRegex, validationTypeMx, validationTypeMxBlacklist, validationTypeSmtp}
}

// Extracts and validates validation type from variadic argument
func variadicValidationType(options []string, defaultValidationType string) (string, error) {
	if len(options) == 0 {
		return defaultValidationType, nil
	}
	validationType := options[0]

	return validationType, validateValidationTypeContext(validationType)
}

// Validates validation type by available values,
// returns error if validation fails
func validateValidationTypeContext(validationType string) error {
	if isIncluded(availableValidationTypes(), validationType) {
		return nil
	}

	return fmt.Errorf(
		"%s is invalid validation type, use one of these: %s",
		validationType,
		availableValidationTypes(),
	)
}

// Returns string by regex pattern capture group index
func regexCaptureGroup(str string, regexPattern string, captureGroup int) string {
	regex, _ := newRegex(regexPattern)

	return regex.FindStringSubmatch(str)[captureGroup]
}

// Returns domain from email string
func emailDomain(email string) string {
	return regexCaptureGroup(email, regexDomainFromEmail, 1)
}

// Returns pointer of copied configuration
func copyConfigurationByPointer(configuration *Configuration) *Configuration {
	config := *configuration
	return &config
}

// Returns a new slice by removing duplicate values in a passed slice
func uniqStrings(strSlice []string) (uniqStrSlice []string) {
	dict := make(map[string]bool)
	for _, item := range strSlice {
		if _, ok := dict[item]; !ok {
			dict[item], uniqStrSlice = true, append(uniqStrSlice, item)
		}
	}

	return uniqStrSlice
}

// Returns a new slice that is a copy of the original slice,
// removing any items that also appear in other slice.
func sliceDiff(slice, otherSlice []string) (diff []string) {
	for _, item := range slice {
		if !isIncluded(otherSlice, item) {
			diff = append(diff, item)
		}
	}

	return diff
}

// Returns server with port number follows {server}:{portNumber} pattern
func serverWithPortNumber(server string, portNumber int) string {
	return fmt.Sprintf("%s:%d", server, portNumber)
}
