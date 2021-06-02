package truemail

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateRegex(t *testing.T) {
	t.Run("Regex validation layer", func(t *testing.T) {
		validatorResult := new(ValidatorResult)
		assert.Equal(t, validateRegex(validatorResult), validatorResult)
	})
}
