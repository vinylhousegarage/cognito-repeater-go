package notfound

// ErrorSimulationResponse is the response returned by the /error/404 endpoint.
// It simulates a 404 Not Found error for testing or demonstration purposes.
type ErrorSimulationResponse struct {
	Message string `json:"message" example:"Simulated 404 Not Found"`
}
