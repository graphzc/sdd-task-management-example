package validator

import (
	"errors"
	"fmt"
	"strings"

	valid "github.com/go-playground/validator/v10"
)

type Validator struct {
	validate *valid.Validate
}

func NewValidator() *Validator {
	return &Validator{
		validate: valid.New(),
	}
}

func (v *Validator) ValidateStruct(payload any) error {
	err := v.validate.Struct(payload)
	if err == nil {
		return nil
	}

	var validationErrs valid.ValidationErrors
	if errors.As(err, &validationErrs) { // Safely unwrap and check for ValidationErrors
		var errMsg strings.Builder
		for _, fieldErr := range validationErrs {
			tmp := strings.Split(fieldErr.StructNamespace(), ".")
			msg := fmt.Sprintf("%s is %s", tmp[len(tmp)-1], fieldErr.Tag())
			msg = strings.ToLower(string(msg[0])) + msg[1:]
			errMsg.WriteString(msg + ", ")
		}

		// Trim trailing comma and space
		finalMsg := strings.TrimSuffix(errMsg.String(), ", ")
		return errors.New(finalMsg)
	}

	// Handle non-validation errors
	return err
}
