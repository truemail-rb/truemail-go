package truemail

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidationMx(t *testing.T) {
	t.Run("when previous validation failed", func(t *testing.T) {
		validator := createValidator(randomEmail(), createConfiguration())
		validatorResult := validator.result

		assert.Equal(t, validatorResult, new(validation).mx(validatorResult))
		assert.Equal(t, usedValidationsByType(ValidationTypeRegex), validator.usedValidations)
	})

	t.Run("DNS validation layer", func(t *testing.T) {
		validator := createValidator(randomEmail(), createConfiguration())
		validatorResult := validator.result
		validatorResult.Success = true

		assert.Equal(t, validatorResult, new(validation).mx(validatorResult))
		assert.Equal(t, usedValidationsByType(ValidationTypeMx), validator.usedValidations)
	})
}
