package whoami

import (
	"encoding/json"
	"net/http"

	"cognito-repeater-go/internal/auth/deps"
	"cognito-repeater-go/internal/httpclient"
)

// @Summary Get user info from access token
// @Description Retrieves user attributes (such as email, sub, etc.) from the Cognito UserInfo endpoint.
// @Description This endpoint expects a Bearer token in the Authorization header and is intended for secure backend use.
// @Description Example: Authorization: Bearer ACCESS_TOKEN_VALUE
// @Tags user
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "User info from Cognito"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /whoami [get]
func NewWhoamiHandler(p deps.WhoamiHandlerProvider, cli httpclient.HTTPClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accessToken, err := ExtractAuthHeaderToken(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		userinfoURL, err := GetUserinfoURL(p.MetadataURL(), cli)
		if err != nil {
			http.Error(w, "failed to get userinfo endpoint", http.StatusInternalServerError)
			return
		}

		userinfo, err := FetchUserinfo(userinfoURL, cli, accessToken)
		if err != nil {
			http.Error(w, "failed to fetch userinfo", http.StatusUnauthorized)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(userinfo)
	}
}
