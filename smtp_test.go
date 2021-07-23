package truemail

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidationSmtpCheck(t *testing.T) {
	t.Run("SMTP validation layer", func(t *testing.T) {
		validator := createValidator(randomEmail(), createConfiguration())
		validatorResult := validator.result

		assert.Equal(t, validatorResult, new(validationSmtp).check(validatorResult))
	})
}
