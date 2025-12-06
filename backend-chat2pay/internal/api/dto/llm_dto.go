package dto

import "chat2pay/internal/entities"

type (
	LLMResponse struct {
		Products []ProductResponse `json:"products"`
		Message  string            `json:"message"`
	}
)

func ToLLM(products *[]entities.Product, message string) LLMResponse {
	if products != nil {
		productResponses := make([]ProductResponse, len(*products))
		for i, product := range *products {
			productResponses[i] = ToProductResponse(&product)
		}

		return LLMResponse{
			Products: productResponses,
			Message:  message,
		}
	}

	return LLMResponse{
		nil,
		message,
	}
}
