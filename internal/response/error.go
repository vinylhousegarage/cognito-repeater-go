package response

import (
    "encoding/json"
    "net/http"

    "go.uber.org/zap"
)

func WriteJSONError(w http.ResponseWriter, status int, msg string, logger *zap.Logger) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)

    if err := json.NewEncoder(w).Encode(ErrorResponse{Error: msg}); err != nil {
        logger.Error("failed to write error response", zap.Error(err))
    }
}
