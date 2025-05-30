package whoami

import (
	"encoding/json"
	"net/http"

	"cognito-repeater-go/internal/auth/deps"
	"cognito-repeater-go/internal/auth/utils"
	"cognito-repeater-go/internal/httpclient"
)

func WhoamiHandler(p deps.WhoamiHandlerProvider, cli httpclient.HTTPClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accessToken, err := utils.ExtractAuthHeaderToken(r)
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
