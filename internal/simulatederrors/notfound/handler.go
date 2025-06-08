package notfound

import (
	"io"
	"log"
	"net/http"

	"go.uber.org/zap"
)

// @Summary Simulate 404 Not Found
// @Description Returns a simulated 404 Not Found error response
// @Tags error
// @Produce plain
// @Failure 404 {string} string "not found 404"
// @Router /error/404 [get]
func NewError404Handler(logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusNotFound)

		if _, err := io.WriteString(w, "not found 404"); err != nil {
			logger.Error("failed to write 404 message", zap.Error(err))
		}
	}
}
