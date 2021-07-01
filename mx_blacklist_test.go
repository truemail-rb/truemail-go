package truemail

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidationMxBlacklist(t *testing.T) {
	t.Run("MX blacklist validation layer", func(t *testing.T) {
		validator := createValidator(randomEmail(), createConfiguration())
		validatorResult := validator.result

		assert.Equal(t, validatorResult, new(validationMxBlacklist).check(validatorResult))
	})
}