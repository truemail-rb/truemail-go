package truemail

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidationRegex(t *testing.T) {
	t.Run("regex validation: successful", func(t *testing.T) {
		validatorResult := createSuccessfulValidatorResult(randomEmail(), createConfiguration())
		new(validationRegex).check(validatorResult)

		assert.True(t, validatorResult.Success)
		assert.Empty(t, validatorResult.Errors)
		assert.Empty(t, validatorResult.usedValidations)
	})

	t.Run("regex validation: failure", func(t *testing.T) {
		validatorResult := createSuccessfulValidatorResult("invalid@email", createConfiguration())
		new(validationRegex).check(validatorResult)

		assert.False(t, validatorResult.Success)
		assert.Equal(t, map[string]string{ValidationTypeRegex: RegexErrorContext}, validatorResult.Errors)
		assert.Empty(t, validatorResult.usedValidations)
	})
}
