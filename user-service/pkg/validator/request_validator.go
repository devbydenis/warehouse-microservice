package pkg

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func Validate(data interface{}) error {
	var errorMessages []string

	err := validate.Struct(data)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Tag() {
				case "required":
					errorMessages = append(errorMessages, fmt.Sprintf("%s is required", err.Field()))
				case "email":
					errorMessages = append(errorMessages, fmt.Sprintf("%s is not a valid email", err.Field()))
				case "min":
					errorMessages = append(errorMessages, fmt.Sprintf("%s must be at least %s characters long", err.Field(), err.Param()))
				// case "max": // <--INI BUAT APA DONGGG??
			}
		}
		return errors.New("Validasi gagal: " + joinMessage(errorMessages))
	}
	return nil
}

func joinMessage(errorMessages []string) string {
	result := ""
	for i, message := range errorMessages {
		if i > 0 {
			result += ", "
		}
		result += message
	}

	return result
}
