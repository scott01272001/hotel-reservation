package types

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())
}

type ValidationError struct {
	message []validator.FieldError
}

func (v ValidationError) Error() string {
	errMsgs := make([]string, len(v.message))
	for i, m := range v.message {
		errMsgs[i] = fmt.Sprintf("Field validation for '%s' failed on the '%s' tag", m.Field(), m.Tag())
	}
	res, err := json.Marshal(errMsgs)
	if err != nil {
		return err.Error()
	}
	return string(res)
}

func ValidateStruct(target interface{}) error {
	targetType := reflect.TypeOf(target)
	if targetType.Kind() == reflect.Ptr {
		targetType = targetType.Elem()
	}
	if targetType.Kind() != reflect.Struct {
		return fmt.Errorf("expected a struct, got %s", targetType.Kind())
	}
	err := validate.Struct(target)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		return ValidationError{
			message: errs,
		}
	}
	return nil
}
