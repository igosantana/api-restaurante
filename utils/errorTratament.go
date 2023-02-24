package utils

import "github.com/go-playground/validator/v10"

type ErrorMsg struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func GetErrorMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "max":
		return "Should be max " + fe.Param()
	case "email":
		return "This field should be email format"
	}
	return "Unknown error"
}
