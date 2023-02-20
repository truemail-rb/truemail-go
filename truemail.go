package truemail

// Validate is main truemail entrypoint. Accepts validation type as option.
// Available types are: regex, mx, mx_blacklist, smtp. By default uses
// validation layer specified in Configuration.ValidationTypeDefault
func Validate(email string, configuration *Configuration, options ...string) (*ValidatorResult, error) {
	validationType, err := variadicValidationType(options, configuration.ValidationTypeDefault)

	if err != nil {
		return nil, err
	}

	return newValidator(email, validationType, configuration).run(), err
}

// IsValid is shortcut for Validate() function. Returns boolean as email validation result.
// Accepts validation type as option. Available types are: regex, mx, mx_blacklist, smtp.
// By default uses validation layer specified in Configuration.ValidationTypeDefault
func IsValid(email string, configuration *Configuration, options ...string) bool {
	validationType, err := variadicValidationType(options, configuration.ValidationTypeDefault)

	if err != nil {
		return false
	}

	return newValidator(email, validationType, configuration).run().Success
}
