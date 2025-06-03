package revoke

import (
	"fmt"
	"net/http"

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

		resp, err := SendRevokeRequest(revokeURL, cli, refreshToken, clientSecret, clientID)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to send revoke request: %v", err), http.StatusInternalServerError)
			return
		}
		defer func() {
			if cerr := resp.Body.Close(); cerr != nil {
				fmt.Printf("failed to close response body: %v\n", cerr)
			}
		}()

		if resp.StatusCode != http.StatusOK {
			http.Error(w, fmt.Sprintf("revocation failed with status: %s", resp.Status), http.StatusBadGateway)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
