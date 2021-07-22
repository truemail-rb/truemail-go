package truemail

import "regexp"

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
