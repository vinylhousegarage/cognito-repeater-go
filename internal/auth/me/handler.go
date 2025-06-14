package me

import (
	"encoding/json"
	"net/http"

	"cognito-repeater-go/internal/auth/deps"
	"cognito-repeater-go/internal/auth/utils"
	"cognito-repeater-go/internal/httpclient"
	"cognito-repeater-go/internal/response"

	"go.uber.org/zap"
)

// @Summary Verify ID token and return user info
// @Description Verifies the provided ID token by validating its signature using Cognito's JWK Set,
// @Description and checking standard claims such as `iss`, `aud`, and `exp`.
// @Description Returns the subject (`sub`) claim from the token if validation succeeds.
// @Description This endpoint accepts a form-encoded POST request and should be called from a secure backend environment.
// @Description
// @Description Example form body:
// @Description   token=ID_TOKEN_VALUE
// @Description   Content-Type: application/x-www-form-urlencoded
// @Tags user
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param token formData string true "ID token to verify"
// @Success 200 {object} me.UserResponse "Returns user sub claim"
// @Failure 400 {object} response.ErrorResponse "Bad Request"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal Server Error"
// @Failure 502 {object} response.ErrorResponse "Internal Bad Gateway"
// @Router /me [post]
func NewMeHandler(
	p deps.MeHandlerProvider,
	c httpclient.HTTPClient,
	logger *zap.Logger,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idToken, err := utils.ExtractFormValue(r, logger)
		if err != nil {
			response.WriteErrorResponse(w, err, logger)
			return
		}

		metadataURL := p.MetadataURL()
		jwksURL, err := GetJWKSURI(metadataURL, c, logger)
		if err != nil {
			response.WriteErrorResponse(w, err, logger)
			return
		}

		jwkSet, err := FetchJWKSet(jwksURL, c, logger)
		if err != nil {
			response.WriteErrorResponse(w, err, logger)
			return
		}

		kid, err := ExtractKIDFromToken(idToken, logger)
		if err != nil {
			response.WriteErrorResponse(w, err, logger)
			return
		}

		jwk, err := FindJWKByKID(kid, jwkSet, logger)
		if err != nil {
			response.WriteErrorResponse(w, err, logger)
			return
		}

		pubKey, err := JWKToRSAPublicKey(jwk, logger)
		if err != nil {
			response.WriteErrorResponse(w, err, logger)
			return
		}

		claims, err := ParseAndVerifyJWT(idToken, pubKey, p.Issuer(), p.Audience())
		if err != nil {
			response.WriteErrorResponse(w, err, logger)
			return
		}

		resp := UserResponse{
			Sub: claims.Subject,
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			logger.Error("failed to write user response", zap.Error(err))
			response.WriteErrorResponse(w, ErrFailedToWriteUserResponse, logger)
			return
		}

		logger.Info("user token verified successfully", zap.String("sub", claims.Subject))
	}
}
