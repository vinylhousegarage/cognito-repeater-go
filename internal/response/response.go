package response

// ErrorResponse is a shared structure for all error responses
type ErrorResponse struct {
    Error string `json:"error" example:"invalid token"`
}
