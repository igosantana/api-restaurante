package response

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

const (
	ErrorTypeFatal        = "Fatal"
	ErrorTypeError        = "Error"
	ErrorTypeValidation   = "Validation Error"
	ErrorTypeInfo         = "Info"
	ErrorTypeDebug        = "Debug"
	ErrorTypeUnauthorized = "Unauthorized"
)

type ErrorDetail struct {
	ErrorType    string `json:"errorType"`
	ErrorMessage string `json:"errorMessage"`
	Field        string `json:"field,omitempty"`
}

type Response struct {
	Data    interface{}   `json:"data,omitempty"`
	Status  int           `json:"status"`
	Error   []ErrorDetail `json:"error,omitempty"`
	Message string        `json:"message"`
}

func GetErrorMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "this field is required"
	case "max":
		return "should be max " + fe.Param()
	case "email":
		return "this field should be email format"
	}
	return "unknown error"
}

func Ok(c *gin.Context, status int, message string, data interface{}) {
	c.JSON(status, Response{
		Data:    data,
		Status:  status,
		Message: message,
	})
}

func BadRequest(c *gin.Context, status int, message string, errors []ErrorDetail) {
	c.JSON(status, Response{
		Error:   errors,
		Status:  status,
		Message: message,
	})
}
