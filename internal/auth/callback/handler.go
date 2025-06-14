package callback

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"cognito-repeater-go/internal/auth/deps"
	"cognito-repeater-go/internal/auth/utils"
	"cognito-repeater-go/internal/httpclient"
	"cognito-repeater-go/internal/response"

	"go.uber.org/zap"
)

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	IDToken      string `json:"id_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

func NewCallbackHandler(
	p deps.CallbackHandlerProvider,
	c httpclient.HTTPClient,
	logger *zap.Logger,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code, err := ValidateCallbackRequest(r)
		if err != nil {
			response.WriteErrorResponse(w, err, logger)
			return
		}

		metadataURL := p.MetadataURL()
		tokenEndpoint, err := GetCallbackURL(metadataURL, c, logger)
		if err != nil {
			response.WriteErrorResponse(w, err, logger)
			return
		}

		bodyStr := BuildTokenRequestBody(code, p)
		req, err := http.NewRequest("POST", tokenEndpoint, strings.NewReader(bodyStr))
		if err != nil {
			response.WriteErrorResponse(w, err, logger)
			return
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Authorization", BuildBasicAuthHeader(p))

		resp, err := c.Do(req)
		if err != nil {
			response.WriteErrorResponse(w, err, logger)
			return
		}
		defer func() {
			if cerr := resp.Body.Close(); cerr != nil {
				logger.Warn("failed to close response body", zap.Error(cerr))
			}
		}()

		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
			response.WriteErrorResponse(w, ErrUnexpectedCallbackStatusCode, logger)
			return
		}

		var tokenResp TokenResponse
		if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
			status := http.StatusInternalServerError
			utils.WritePlainError(w, status, err, logger)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		logger.Info("token successfully issued",
			zap.String("client_id", p.UserPoolClientIDValue()),
			zap.String("token_type", tokenResp.TokenType),
		)

		if err := json.NewEncoder(w).Encode(tokenResp); err != nil {
			status := http.StatusInternalServerError
			utils.WritePlainError(w, status, err, logger)
			return
		}
	}
}
