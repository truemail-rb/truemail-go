package truemail

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidationMxBlacklistCheck(t *testing.T) {
	t.Run("MX blacklist validation: successful", func(t *testing.T) {
		validatorResult := createSuccessfulValidatorResult(randomEmail(), createConfiguration())
		new(validationMxBlacklist).check(validatorResult)

		assert.True(t, validatorResult.Success)
		assert.Empty(t, validatorResult.Errors)
		assert.Empty(t, validatorResult.usedValidations)
	})

	t.Run("MX blacklist validation: failure", func(t *testing.T) {
		blacklistedMxIpAddress := randomIpAddress()
		configuration, _ := NewConfiguration(ConfigurationAttr{VerifierEmail: randomEmail(), BlacklistedMxIpAddresses: []string{blacklistedMxIpAddress}})
		validatorResult := createSuccessfulValidatorResult(randomEmail(), configuration)
		validatorResult.MailServers = []string{blacklistedMxIpAddress}
		new(validationMxBlacklist).check(validatorResult)

		assert.False(t, validatorResult.Success)
		assert.Equal(t, map[string]string{validationTypeMxBlacklist: mxBlacklistErrorContext}, validatorResult.Errors)
		assert.Empty(t, validatorResult.usedValidations)
	})
}
