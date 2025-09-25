package middleware

import (
	"encoding/json"
	"errors"

	"github.com/go-playground/validator"
)

type CustomValidator struct {
	Validator *validator.Validate
}

type ValidationErrors struct {
	Errors map[string][]string `json:"errors"`
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.Validator.Struct(i); err != nil {

		var invalidError *validator.InvalidValidationError
		if errors.As(err, &invalidError) {
			return err
		}

		var validateErrs validator.ValidationErrors
		if errors.As(err, &validateErrs) {
			validationErrors := make(map[string][]string)
			for _, err3 := range validateErrs {
				fieldName := err3.Field()
				validationErrors["S"+fieldName] = append(validationErrors[fieldName], err3.Tag())
			}

			if errJSON, err := json.Marshal(ValidationErrors{Errors: validationErrors}); err == nil {
				return errors.New(string(errJSON))
			}
		}
		// Optionally, you could return the error to give each route more control over the status code

		return err
	}
	return nil
}
