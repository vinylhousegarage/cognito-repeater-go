package whoami

import (
	"encoding/json"
	"errors"
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
			switch {
			case errors.Is(err, ErrMissingAuthorizationHeader):
				logger.Warn("missing authorization header", zap.Error(err))
				response.WriteJSONError(w, http.StatusUnauthorized, err.Error(), logger)
			case errors.Is(err, ErrInvalidAuthorizationHeaderFormat):
				logger.Warn("invalid authorization header format", zap.Error(err))
				response.WriteJSONError(w, http.StatusBadRequest, err.Error(), logger)
			default:
				status := http.StatusInternalServerError
				logger.Error("unexpected error while extracting token", zap.Error(err))
				response.WriteJSONError(w, status, http.StatusText(status), logger)
			}
			return
		}

		userinfoURL, err := GetUserinfoURL(p.MetadataURL(), cli, logger)
		if err != nil {
			switch {
			case errors.Is(err, ErrFailedToCreateRequest):
				response.WriteJSONError(w, http.StatusInternalServerError, err.Error(), logger)
			case errors.Is(err, ErrFailedToFetchMetadata):
				response.WriteJSONError(w, http.StatusBadGateway, err.Error(), logger)
			case errors.Is(err, ErrFailedToReadMetadataResponse):
				response.WriteJSONError(w, http.StatusBadGateway, err.Error(), logger)
			case errors.Is(err, ErrUnexpectedMetadataStatusCode):
				response.WriteJSONError(w, http.StatusBadGateway, err.Error(), logger)
			case errors.Is(err, ErrFailedToDecodeMetadata):
				response.WriteJSONError(w, http.StatusBadGateway, err.Error(), logger)
			case errors.Is(err, ErrMissingUserinfoEndpoint):
				response.WriteJSONError(w, http.StatusBadGateway, err.Error(), logger)
			default:
				logger.Error("unexpected error while extracting token", zap.Error(err))
				response.WriteJSONError(w, http.StatusInternalServerError, "internal server error", logger)
			}
			return
		}

		userinfo, err := FetchUserinfo(userinfoURL, cli, accessToken, logger)
		if err != nil {
			switch {
			case errors.Is(err, ErrFailedToCreateUserinfoRequest):
				response.WriteJSONError(w, http.StatusInternalServerError, err.Error(), logger)
			case errors.Is(err, ErrFailedToFetchUserinfoRequest):
				response.WriteJSONError(w, http.StatusBadGateway, err.Error(), logger)
			case errors.Is(err, ErrFailedToReadUserinfoResponse):
				response.WriteJSONError(w, http.StatusBadGateway, err.Error(), logger)
			case errors.Is(err, ErrUnexpectedUserinfoStatusCode):
				response.WriteJSONError(w, http.StatusBadGateway, err.Error(), logger)
			default:
				logger.Error("unexpected error while extracting token", zap.Error(err))
				response.WriteJSONError(w, http.StatusInternalServerError, "internal server error", logger)
			}
			return
		}

		subRaw, ok := userinfo["sub"]
		if !ok {
			logger.Warn("missing subject (sub)", zap.Any("userinfo", userinfo))
			response.WriteJSONError(w, http.StatusUnauthorized, ErrMissingSubject.Error(), logger)
			return
		}
		sub, ok := subRaw.(string)
		if !ok {
			logger.Warn("subject (sub) claim is not a string", zap.Any("subRaw", subRaw))
			response.WriteJSONError(w, http.StatusUnauthorized, ErrSubjectIsNotString.Error(), logger)
			return
		}

		resp := UserInfoResponse{
			Sub: sub,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			logger.Error("failed to write response", zap.Any("resp", resp), zap.Error(err))
		}
	}
}
