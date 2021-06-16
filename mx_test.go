package truemail

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateMx(t *testing.T) {
	t.Run("DNS validation layer", func(t *testing.T) {
		validatorResult := createValidatorResult(createRandomEmail(), createConfiguration())
		assert.Equal(t, validateMx(validatorResult), validatorResult)
	})
}
