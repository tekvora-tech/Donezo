package validator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Validator wrapper go-playground/validator
type Validator struct {
	validate *validator.Validate
}

// ValidationError error validasi per field
type ValidationError struct {
	Field   string `json:"field"`   // Nama field: "email"
	Message string `json:"message"` // Pesan: "email harus valid"
}

// ValidationErrors kumpulan error validasi
type ValidationErrors []ValidationError

// Error implementasi error interface
func (ve ValidationErrors) Error() string {
	var messages []string
	for _, e := range ve {
		messages = append(messages, fmt.Sprintf("%s: %s", e.Field, e.Message))
	}
	return strings.Join(messages, "; ")
}

// New buat instance validator
func New() *Validator {
	v := validator.New()

	// Custom tag name (pakai json tag untuk nama field)
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return fld.Name
		}
		return name
	})

	return &Validator{validate: v}
}

// Validate cek struct
// Return ValidationErrors (bisa di-cast) atau error biasa
func (v *Validator) Validate(s interface{}) error {
	if err := v.validate.Struct(s); err != nil {
		// Cast ke validator.ValidationErrors
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			return v.parseErrors(validationErrors)
		}
		return err
	}
	return nil
}

// ValidateVar validasi single variable
func (v *Validator) ValidateVar(field interface{}, tag string) error {
	return v.validate.Var(field, tag)
}

// parseErrors convert validator errors ke format kita
func (v *Validator) parseErrors(errors validator.ValidationErrors) ValidationErrors {
	var result ValidationErrors

	for _, err := range errors {
		result = append(result, ValidationError{
			Field:   err.Field(),
			Message: v.translateError(err),
		})
	}

	return result
}

// translateError ubah tag ke pesan yang readable
func (v *Validator) translateError(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s wajib diisi", err.Field())
	case "email":
		return fmt.Sprintf("%s harus format email yang valid", err.Field())
	case "min":
		return fmt.Sprintf("%s minimal %s karakter", err.Field(), err.Param())
	case "max":
		return fmt.Sprintf("%s maksimal %s karakter", err.Field(), err.Param())
	case "gte":
		return fmt.Sprintf("%s harus lebih besar atau sama dengan %s", err.Field(), err.Param())
	case "lte":
		return fmt.Sprintf("%s harus lebih kecil atau sama dengan %s", err.Field(), err.Param())
	case "oneof":
		return fmt.Sprintf("%s harus salah satu dari: %s", err.Field(), err.Param())
	case "numeric":
		return fmt.Sprintf("%s harus angka", err.Field())
	case "alphanum":
		return fmt.Sprintf("%s hanya boleh huruf dan angka", err.Field())
	case "uuid":
		return fmt.Sprintf("%s harus format UUID", err.Field())
	case "url":
		return fmt.Sprintf("%s harus format URL", err.Field())
	case "datetime":
		return fmt.Sprintf("%s harus format datetime: %s", err.Field(), err.Param())
	default:
		return fmt.Sprintf("%s tidak valid", err.Field())
	}
}