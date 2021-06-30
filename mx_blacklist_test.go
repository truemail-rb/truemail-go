package truemail

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidationMxBlacklist(t *testing.T) {
	t.Run("when previous validation failed", func(t *testing.T) {
		validator := createValidator(randomEmail(), createConfiguration())
		validatorResult := validator.result

		assert.Equal(t, validatorResult, new(validation).mxBlacklist(validatorResult))
		assert.Equal(t, usedValidationsByType(ValidationTypeRegex), validator.usedValidations)
	})

	t.Run("MX blacklist validation layer", func(t *testing.T) {
		validator := createValidator(randomEmail(), createConfiguration())
		validatorResult := validator.result
		validatorResult.Success = true

		assert.Equal(t, validatorResult, new(validation).mxBlacklist(validatorResult))
		assert.Equal(t, usedValidationsByType(ValidationTypeMx), validator.usedValidations)
	})
}
