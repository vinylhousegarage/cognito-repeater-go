package revoke

import (
	"fmt"
	"net/http"
	"strings"

	"cognito-repeater-go/internal/auth/deps"
	"cognito-repeater-go/internal/auth/utils"
	"cognito-repeater-go/internal/httpclient"
)

func NewRevokeHandler(p deps.RevokeHandlerProvider, cli httpclient.HTTPClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		refreshToken, err := utils.ExtractFormValue(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		revokeURL, err := GetRevokeURL(p.MetadataURL(), cli)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to get revoke endpoint: %v", err), http.StatusInternalServerError)
			return
		}

		clientSecret := p.ClientSecretValue()
		clientID := p.UserPoolClientIDValue()

		clientSecret = strings.TrimSpace(clientSecret)
		clientID = strings.TrimSpace(clientID)

		resp, err := SendRevokeRequest(revokeURL, cli, refreshToken, clientID, clientSecret)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer func() {
			if err := resp.Body.Close(); err != nil {
				fmt.Printf("failed to close response body: %v\n", err)
			}
		}()

		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
			http.Error(w, fmt.Sprintf("revocation failed with status: %s", resp.Status), http.StatusBadGateway)
			return
		}

		fmt.Printf("revocation succeeded with status: %s\n", resp.Status)
		w.WriteHeader(http.StatusNoContent)
	}
}
