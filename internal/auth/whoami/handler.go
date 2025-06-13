package whoami

import (
	"encoding/json"
	"net/http"

	"cognito-repeater-go/internal/auth/deps"
	"cognito-repeater-go/internal/httpclient"
	"cognito-repeater-go/internal/response"

	"go.uber.org/zap"
)

// @Summary Get user info from access token
// @Description Retrieves the user sub from the Cognito UserInfo endpoint.
// @Description This endpoint expects a Bearer token in the Authorization header.
// @Description Example: Authorization: Bearer ACCESS_TOKEN_VALUE
// @Tags user
// @Produce json
// @Security BearerAuth
// @Success 200 {object} whoami.UserInfoResponse "User info from Cognito"
// @Failure 400 {object} response.ErrorResponse "Bad Request"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal Server Error"
// @Failure 502 {object} response.ErrorResponse "Bad Gateway"
// @Router /whoami [get]
func NewWhoamiHandler(
	p deps.WhoamiHandlerProvider,
	cli httpclient.HTTPClient,
	logger *zap.Logger,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accessToken, err := ExtractAuthHeaderToken(r)
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

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			logger.Error("failed to write response", zap.Any("resp", resp), zap.Error(err))
		}
	}
}
