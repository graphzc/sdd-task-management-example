package dto

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type HealthCheckResponse struct {
	Status string `json:"status"`
}

type MessageResponse struct {
	Message string `json:"message"`
}
