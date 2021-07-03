package truemail

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidationMx(t *testing.T) {
	t.Run("DNS validation layer", func(t *testing.T) {
		validator := createValidator(randomEmail(), createConfiguration())
		validatorResult := validator.result

		assert.Equal(t, validatorResult, new(validationMx).check(validatorResult))
	})
}
