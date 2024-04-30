package validation

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"strings"
)

func ValidationError(errs validator.ValidationErrors) error {
	var errMsgs []string
	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is required", err.Field()))
		case "url":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s must be a valid url", err.Field()))
		case "string":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s must be a string", err.Field()))
		case "float64":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s must be a float", err.Field()))
		case "int":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s must be a int", err.Field()))
		case "min":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s must be a min %s", err.Field(), err.Param()))
		case "max":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s must be a max %s", err.Field(), err.Param()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is invalid", err.Field()))
		}
	}

	return fmt.Errorf(strings.Join(errMsgs, ", "))
}
