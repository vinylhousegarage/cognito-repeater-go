package health

// HealthResponse represents a simple health check response.
// Typically used for /health endpoint.
type HealthResponse struct {
	Message string `json:"message" example:"pong"`
}
