package presenter

import "github.com/gofiber/fiber/v2"

type Response struct {
	Code   int         `json:"code"`
	Data   interface{} `json:"data,omitempty"`
	Errors error       `json:"errors,omitempty"`
}

// SuccessResponseSwagger untuk dokumentasi Swagger
type SuccessResponseSwagger struct {
	Status bool        `json:"status" example:"true"`
	Data   interface{} `json:"data"`
	Error  *string     `json:"error"`
}

// ErrorResponseSwagger untuk dokumentasi Swagger
type ErrorResponseSwagger struct {
	Status bool   `json:"status" example:"false"`
	Data   string `json:"data" example:""`
	Error  string `json:"error" example:"error message"`
}

func (r *Response) WithCode(code int) *Response {
	r.Code = code
	return r
}

func (r *Response) WithData(data interface{}) *Response {
	r.Data = data
	return r
}

func (r *Response) WithError(err error) *Response {
	r.Errors = err
	return r
}

func ErrorResponse(err error) *fiber.Map {
	return &fiber.Map{
		"status": false,
		"data":   "",
		"error":  err.Error(),
	}
}

func SuccessResponse(data interface{}) *fiber.Map {
	return &fiber.Map{
		"status": true,
		"data":   data,
		"error":  nil,
	}
}

func CreatedResponse() *fiber.Map {
	return &fiber.Map{
		"status": true,
		"data":   "created",
		"error":  nil,
	}
}
