package truemail

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidationRegex(t *testing.T) {
	t.Run("Regex validation layer", func(t *testing.T) {
		validator := createValidator(randomEmail(), createConfiguration())
		validatorResult := validator.result

		assert.Equal(t, validatorResult, new(validationRegex).check(validatorResult))
	})
}
