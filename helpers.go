package truemail

import (
	"fmt"
	"regexp"
)

// package helpers functions

func isIncluded(slice []string, target string) bool {
	for _, item := range slice {
		if item == target {
			return true
		}
	}

	return false
}

func isIntersected(baseSlice []string, targetSlice []string) bool {
	for _, item := range targetSlice {
		if isIncluded(baseSlice, item) {
			return true
		}
	}

	return false
}

func newRegex(regexPattern string) (*regexp.Regexp, error) {
	return regexp.Compile(regexPattern)
}

func matchRegex(strContext, regexPattern string) bool {
	regex, err := newRegex(regexPattern)
	if err != nil {
		return false
	}

	return regex.MatchString(strContext)
}

func availableValidationTypes() []string {
	return []string{ValidationTypeRegex, ValidationTypeMx, ValidationTypeMxBlacklist, ValidationTypeSMTP}
}

func variadicValidationType(options []string, defaultValidationType string) (string, error) {
	if len(options) == 0 {
		return defaultValidationType, nil
	}
	validationType := options[0]

	return validationType, validateValidationTypeContext(validationType)
}

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

func regexCaptureGroup(str string, regexPattern string, captureGroup int) string {
	regex, _ := newRegex(regexPattern)

	return regex.FindStringSubmatch(str)[captureGroup]
}

func emailDomain(email string) string {
	return regexCaptureGroup(email, RegexDomainFromEmail, 1)
}

func copyConfigurationByPointer(configuration *configuration) *configuration {
	config := *configuration
	return &config
}

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
