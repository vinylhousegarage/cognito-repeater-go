package revoke

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"cognito-repeater-go/internal/auth/deps"
	"cognito-repeater-go/internal/auth/utils"
	"cognito-repeater-go/internal/httpclient"
	"cognito-repeater-go/internal/response"
)

func writeJSONError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(response.ErrorResponse{
		Error: msg,
	}); err != nil {
		log.Printf("failed to write error response: %v", err)
	}
}

// @Summary Revoke refresh token
// @Description Revokes a refresh token by calling the Cognito revocation endpoint.
// @Description This endpoint expects a form-encoded POST request and should be called from a secure backend environment.
// @Description To revoke a refresh token, send a POST request to this endpoint with the following body:
// @Description token=REFRESH_TOKEN_VALUE
// @Description Content-Type: application/x-www-form-urlencoded
// @Tags auth
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param token formData string true "Refresh token to be revoked"
// @Success 204 {string} string "No Content"
// @Failure 400 {object} response.ErrorResponse "Bad Request"
// @Failure 502 {object} response.ErrorResponse "Bad Gateway"
// @Failure 500 {object} response.ErrorResponse "Internal Server Error"
// @Router /revoke [post]
func NewRevokeHandler(p deps.RevokeHandlerProvider, cli httpclient.HTTPClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		refreshToken, err := utils.ExtractFormValue(r)
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}

		revokeURL, err := GetRevokeURL(p.MetadataURL(), cli)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "failed to resolve revocation endpoint")
			return
		}

		clientSecret := p.ClientSecretValue()
		clientID := p.UserPoolClientIDValue()

		clientSecret = strings.TrimSpace(clientSecret)
		clientID = strings.TrimSpace(clientID)

		resp, err := SendRevokeRequest(revokeURL, cli, refreshToken, clientID, clientSecret)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "failed to get userinfo endpoint")
			return
		}
		defer func() {
			if err := resp.Body.Close(); err != nil {
				log.Printf("failed to close response body: %v", err)
			}
		}()

		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
			writeJSONError(w, http.StatusBadGateway, fmt.Sprintf("revocation failed with status: %s", resp.Status))
			return
		}

		log.Printf("revocation succeeded with status: %s", resp.Status)
		w.WriteHeader(http.StatusNoContent)
	}
}
