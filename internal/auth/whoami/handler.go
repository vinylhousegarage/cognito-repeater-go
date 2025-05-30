package whoami

import (
	"encoding/json"
	"net/http"

	"cognito-repeater-go/internal/auth/deps"
	"cognito-repeater-go/internal/auth/utils"
)

func WhoamiHandler(d deps.HandlerDependencies) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accessToken, err := utils.ExtractAuthHeaderToken(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		userinfoURL, err := GetUserinfoURL(d.Config.MetadataURL(), d.HTTPClient)
		if err != nil {
			http.Error(w, "failed to get userinfo endpoint", http.StatusInternalServerError)
			return
		}

		userinfo, err := FetchUserinfo(d.HTTPClient, userinfoURL, accessToken)
		if err != nil {
			http.Error(w, "failed to fetch userinfo", http.StatusUnauthorized)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(userinfo)
	}
}
