package types

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())
}

func ValidateStruct(target interface{}) (map[string]string, error) {
	errors := map[string]string{}

	targetType := reflect.TypeOf(target)
	if targetType.Kind() == reflect.Ptr {
		targetType = targetType.Elem()
	}
	if targetType.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected a struct, got %s", targetType.Kind())
	}
	if err := validate.Struct(target); err != nil {
		var builder strings.Builder

		for _, e := range err.(validator.ValidationErrors) {
			builder.Reset()

			builder.WriteString(fmt.Sprintf("constrain=%s", e.Tag()))
			if e.Param() != "" {
				builder.WriteString(fmt.Sprintf(", param=%s", e.Param()))
			}

			errors[e.Field()] = builder.String()
		}
		return errors, nil
	}
	return nil, nil
}
