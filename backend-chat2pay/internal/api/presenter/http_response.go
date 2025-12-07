package presenter

import "github.com/gofiber/fiber/v2"

type Response struct {
	Code   int         `json:"-"`
	Status bool        `json:"status"`
	Data   interface{} `json:"data,omitempty"`
	Error  string      `json:"error,omitempty"`
	Errors error       `json:"-"` // internal use only, not serialized
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

func NewResponse() *Response {
	return &Response{}
}

func (r *Response) WithCode(code int) *Response {
	r.Code = code
	r.Status = code >= 200 && code < 300
	return r
}

func (r *Response) WithData(data interface{}) *Response {
	r.Data = data
	return r
}

func (r *Response) WithError(err error) *Response {
	r.Errors = err
	if err != nil {
		r.Error = err.Error()
	}
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
