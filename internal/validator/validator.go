package validator

import (
	"regexp"
	"slices"
	"strings"
	"unicode/utf8"
)

type Validator struct {
	NonFieldErrors []string
	FieldErrors    map[string]string
}

var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0 && len(v.NonFieldErrors) == 0
}

func (v *Validator) AddFieldError(field, message string) {
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	if _, exists := v.FieldErrors[field]; !exists {
		v.FieldErrors[field] = message
	}
}

func (v *Validator) AddNonFieldError(message string) {
	v.NonFieldErrors = append(v.NonFieldErrors, message)
}
func (v *Validator) CheckField(ok bool, field, message string) {
	if !ok {
		v.AddFieldError(field, message)
	}
}
func MaxChars(field string, chars int) bool {
	return utf8.RuneCountInString(field) <= chars
}

func MinChars(value string, n int) bool {
	return utf8.RuneCountInString(value) >= n
}

func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}
func NotBlank(field string) bool {
	return strings.TrimSpace(field) != ""
}

func PermittedValue[T comparable](value T, permittedValues ...T) bool {
	return slices.Contains(permittedValues, value)
}
