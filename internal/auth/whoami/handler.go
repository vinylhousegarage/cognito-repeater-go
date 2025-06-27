package whoami

import (
	"encoding/json"
	"net/http"

	"cognito-repeater-go/internal/auth/deps"
	"cognito-repeater-go/internal/auth/utils"
	"cognito-repeater-go/internal/httpclient"
	"cognito-repeater-go/internal/response"

	"go.uber.org/zap"
)

func NewWhoamiHandler(
	p deps.WhoamiHandlerProvider,
	cli httpclient.HTTPClient,
	logger *zap.Logger,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Vary", "Origin")
			w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Authorization")
			w.WriteHeader(http.StatusNoContent)
			return
		}

		accessToken, err := utils.ExtractAuthHeaderToken(r)
		if err != nil {
			response.WriteErrorResponse(w, err, logger)
			return
		}

		userinfoURL, err := GetUserinfoURL(p.MetadataURL(), cli, logger)
		if err != nil {
			response.WriteErrorResponse(w, err, logger)
			return
		}

		userinfo, err := FetchUserinfo(userinfoURL, cli, accessToken, logger)
		if err != nil {
			response.WriteErrorResponse(w, err, logger)
			return
		}

		if userinfo.Sub == "" {
			logger.Warn("missing subject (sub)", zap.Any("userinfo", userinfo))
			response.WriteErrorResponse(w, ErrMissingSubject, logger)
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Vary", "Origin")
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(userinfo); err != nil {
			logger.Error("failed to write response", zap.Any("resp", userinfo), zap.Error(err))
		}
	}
}
