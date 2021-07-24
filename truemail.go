package truemail

// Main truemail function. Accepts validation type as option. Available types are: regex, mx,
// mx_blacklist, smtp. By default uses validation layer specified in configuration.validationTypeDefault
func Validate(email string, configuration *configuration, options ...string) (*validatorResult, error) {
	validationType, err := variadicValidationType(options, configuration.ValidationTypeDefault)

	if err != nil {
		return nil, err
	}

	return newValidator(email, validationType, configuration).run(), err
}

// Shortcut for Validate() function. Returns boolean as email validation result.
// Accepts validation type as option. Available types are: regex, mx, mx_blacklist, smtp.
// By default uses validation layer specified in configuration.validationTypeDefault
func IsValid(email string, configuration *configuration, options ...string) bool {
	validationType, err := variadicValidationType(options, configuration.ValidationTypeDefault)

	if err != nil {
		return false
	}

	return newValidator(email, validationType, configuration).run().Success
}
