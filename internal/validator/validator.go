package validator

import (
	"slices"
	"strings"
	"unicode/utf8"
)

type Validator struct {
	Errors map[string]string
}

func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

func (v *Validator) AddFieldError(field, message string) {
	if v.Errors == nil {
		v.Errors = make(map[string]string)
	}

	if _, exists := v.Errors[field]; !exists {
		v.Errors[field] = message
	}
}
func (v *Validator) CheckField(ok bool, field, message string) {
	if !ok {
		v.AddFieldError(field, message)
	}
}
func MaxChars(field string, chars int) bool {
	return utf8.RuneCountInString(field) <= chars
}

func NotBlank(field string) bool {
	return strings.TrimSpace(field) != ""
}

func PermittedValue[T comparable](value T, permittedValues ...T) bool {
	return slices.Contains(permittedValues, value)
}
