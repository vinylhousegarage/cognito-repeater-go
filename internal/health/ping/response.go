package ping

// PingResponse represents a simple health check response.
// Typically used for /ping endpoint.
type PingResponse struct {
    Message string `json:"message" example:"pong"`
}
