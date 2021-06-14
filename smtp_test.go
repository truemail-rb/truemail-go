package truemail

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateSMTP(t *testing.T) {
	t.Run("SMTP validation layer", func(t *testing.T) {
		validatorResult := new(validatorResult)
		assert.Equal(t, validateSMTP(validatorResult), validatorResult)
	})
}
