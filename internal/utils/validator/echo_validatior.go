package validator

import (
	"github.com/graphzc/sdd-task-management-example/internal/utils/servererr"
)

func (v *Validator) Validate(payload any) error {
	err := v.ValidateStruct(payload)
	if err != nil {
		return servererr.NewError(
			servererr.ErrorCodeBadRequest,
			err.Error(),
		)
	}

	return nil
}
