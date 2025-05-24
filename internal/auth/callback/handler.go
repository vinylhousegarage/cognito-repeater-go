package callback

import (
	"cognito-repeater-go/internal/auth/deps"
	"encoding/json"
	"net/http"
	"strings"
)

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	IDToken      string `json:"id_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

func CallbackHandler(d deps.HandlerDependencies) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code, err := ValidateCallbackRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		tokenEndpoint, err := GetCallbackURL(d.HTTPClient, d.Config)
		if err != nil {
			http.Error(w, "failed to get token endpoint", http.StatusInternalServerError)
			return
		}

		bodyStr := BuildTokenRequestBody(code, d.Config)
		req, err := http.NewRequest("POST", tokenEndpoint, strings.NewReader(bodyStr))
		if err != nil {
			http.Error(w, "failed to create request", http.StatusInternalServerError)
			return
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Authorization", BuildBasicAuthHeader(d.Config))

		resp, err := d.HTTPClient.Do(req)
		if err != nil {
			http.Error(w, "failed to send token request", http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			http.Error(w, "token endpoint returned error", http.StatusBadGateway)
			return
		}

		var tokenResp TokenResponse
		if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
			http.Error(w, "failed to decode token response", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(tokenResp); err != nil {
			http.Error(w, "failed to encode token response", http.StatusInternalServerError)
		}
	}
}
