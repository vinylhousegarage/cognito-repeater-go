package whoami

import (
	"encoding/json"
	"log"
	"net/http"

	"cognito-repeater-go/internal/auth/deps"
	"cognito-repeater-go/internal/httpclient"
	"cognito-repeater-go/internal/response"

	"go.uber.org/zap"
)

func writeJSONError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(response.ErrorResponse{
		Error: msg,
	})
	if err != nil {
		log.Printf("failed to write error response: %v", err)
	}
}

// @Summary Get user info from access token
// @Description Retrieves the user sub from the Cognito UserInfo endpoint.
// @Description This endpoint expects a Bearer token in the Authorization header.
// @Description Example: Authorization: Bearer ACCESS_TOKEN_VALUE
// @Tags user
// @Produce json
// @Security BearerAuth
// @Success 200 {object} whoami.UserInfoResponse "User info from Cognito"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal Server Error"
// @Router /whoami [get]
func NewWhoamiHandler(
	p deps.WhoamiHandlerProvider,
	cli httpclient.HTTPClient,
	logger *zap.Logger,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_ = logger

		accessToken, err := ExtractAuthHeaderToken(r)
		if err != nil {
			writeJSONError(w, http.StatusUnauthorized, err.Error())
			return
		}

		userinfoURL, err := GetUserinfoURL(p.MetadataURL(), cli)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "failed to get userinfo endpoint")
			return
		}

		userinfo, err := FetchUserinfo(userinfoURL, cli, accessToken)
		if err != nil {
			writeJSONError(w, http.StatusUnauthorized, "failed to fetch userinfo")
			return
		}

		subRaw, ok := userinfo["sub"]
		if !ok {
			writeJSONError(w, http.StatusUnauthorized, `"sub" claim is missing`)
			return
		}
		sub, ok := subRaw.(string)
		if !ok {
			writeJSONError(w, http.StatusUnauthorized, `"sub" claim is not a string`)
			return
		}

		resp := UserInfoResponse{
			Sub: sub,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("failed to write response: %v", err)
		}
	}
}
