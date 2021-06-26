package truemail

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateMxBlacklist(t *testing.T) {
	t.Run("MX blacklist validation layer", func(t *testing.T) {
		validatorResult := createValidatorResult(randomEmail(), createConfiguration())
		assert.Equal(t, validateMxBlacklist(validatorResult), validatorResult)
	})
}
