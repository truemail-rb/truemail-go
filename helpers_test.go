package truemail

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsIncluded(t *testing.T) {
	var item string

	t.Run("item found in slice", func(t *testing.T) {
		assert.True(t, isIncluded([]string{item}, item))
	})

	t.Run("item not found in slice", func(t *testing.T) {
		assert.False(t, isIncluded([]string{}, item))
	})
}

func TestIsIntersected(t *testing.T) {
	var item string

	t.Run("item from target slice found in base slice", func(t *testing.T) {
		assert.True(t, isIntersected([]string{item}, []string{item}))
	})

	t.Run("item from target slice not found in base slice", func(t *testing.T) {
		assert.False(t, isIntersected([]string{}, []string{item}))
	})
}

func TestNewRegex(t *testing.T) {
	t.Run("valid regex pattern", func(t *testing.T) {
		regexPattern := EmptyString
		actualRegex, err := newRegex(regexPattern)
		expectedRegex, _ := regexp.Compile(regexPattern)

		assert.Equal(t, expectedRegex, actualRegex)
		assert.NoError(t, err)
	})

	t.Run("invalid regex pattern", func(t *testing.T) {
		actualRegex, err := newRegex(`\K`)

		assert.Nil(t, actualRegex)
		assert.Error(t, err)
	})
}

func TestMatchRegex(t *testing.T) {
	t.Run("valid regex pattern, matched string", func(t *testing.T) {
		assert.True(t, matchRegex(EmptyString, EmptyString))
	})

	t.Run("valid regex pattern, not matched string", func(t *testing.T) {
		assert.False(t, matchRegex("42", `\D+`))
	})

	t.Run("invalid regex pattern", func(t *testing.T) {
		assert.False(t, matchRegex(EmptyString, `\K`))
	})
}

func TestAvailableValidationTypes(t *testing.T) {
	t.Run("slice of available validation types", func(t *testing.T) {
		assert.Equal(t, []string{"regex", "mx", "mx_blacklist", "smtp"}, availableValidationTypes())
	})
}

func TestVariadicValidationType(t *testing.T) {
	t.Run("without validation type", func(t *testing.T) {
		result, err := variadicValidationType([]string{}, ValidationTypeMx)

		assert.NoError(t, err)
		assert.Equal(t, ValidationTypeMx, result)
	})

	t.Run("valid validation type", func(t *testing.T) {
		validationType := ValidationTypeRegex
		result, err := variadicValidationType([]string{validationType}, ValidationTypeMx)

		assert.NoError(t, err)
		assert.Equal(t, validationType, result)
	})

	t.Run("invalid validation type", func(t *testing.T) {
		invalidValidationType := "invalid type"
		result, err := variadicValidationType([]string{invalidValidationType}, ValidationTypeMx)
		errorMessage := fmt.Sprintf("%s is invalid validation type, use one of these: [regex mx mx_blacklist smtp]", invalidValidationType)

		assert.EqualError(t, err, errorMessage)
		assert.Equal(t, invalidValidationType, result)
	})
}

func TestValidateValidationTypeContext(t *testing.T) {
	for _, validValidationType := range []string{ValidationTypeRegex, ValidationTypeMx, ValidationTypeSMTP} {
		t.Run("valid validation type", func(t *testing.T) {
			assert.NoError(t, validateValidationTypeContext(validValidationType))
		})
	}

	t.Run("invalid validation type", func(t *testing.T) {
		invalidType := "invalid type"
		errorMessage := fmt.Sprintf("%s is invalid validation type, use one of these: [regex mx mx_blacklist smtp]", invalidType)

		assert.EqualError(t, validateValidationTypeContext(invalidType), errorMessage)
	})
}
