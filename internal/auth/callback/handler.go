package callback

import (
	"encoding/json"
	"net/http"
	"strings"

	"cognito-repeater-go/internal/config"
	"cognito-repeater-go/internal/httpclient"
)

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	IDToken      string `json:"id_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

type CallbackHandlerDependencies struct {
	Config      *config.Config
	HTTPClient  httpclient.HTTPClient
	URLProvider CallbackURLProvider
}

func CallbackHandler(deps CallbackHandlerDependencies) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code, err := ValidateCallbackRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		tokenEndpoint, err := deps.URLProvider.GetCallbackURL(deps.Config)
		if err != nil {
			http.Error(w, "failed to get token endpoint", http.StatusInternalServerError)
			return
		}

		bodyStr := BuildTokenRequestBody(code, deps.Config)
		req, err := http.NewRequest("POST", tokenEndpoint, strings.NewReader(bodyStr))
		if err != nil {
			http.Error(w, "failed to create request", http.StatusInternalServerError)
			return
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Authorization", BuildBasicAuthHeader(deps.Config))

		resp, err := deps.HTTPClient.Do(req)
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
		json.NewEncoder(w).Encode(tokenResp)
	}
}
