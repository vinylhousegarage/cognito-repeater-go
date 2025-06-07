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

func writeJSONError(w http.ResponseWriter, status int, msg string, logger *zap.Logger) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(response.ErrorResponse{
		Error: msg,
	})
	if err != nil {
		logger.Error("failed to write error response", zap.Error(err))
	}
}

// @Summary Verify ID token and return user info
// @Description Verifies the provided ID token using Cognito's JWKS and returns the subject (sub) claim.
// @Description This endpoint expects a form-encoded POST request and should be called from a secure backend environment.
// @Description To verify an ID token, send a POST request with the following body:
// @Description token=ID_TOKEN_VALUE
// @Description Content-Type: application/x-www-form-urlencoded
// @Tags user
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param token formData string true "ID token to verify"
// @Success 200 {object} me.UserResponse "Returns user sub claim"
// @Failure 400 {object} response.ErrorResponse "Bad Request"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal Server Error"
// @Router /me [post]
func NewMeHandler(
	p deps.MeHandlerProvider,
	c httpclient.HTTPClient,
	logger *zap.Logger,
	) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenStr, err := utils.ExtractFormValue(r)
		if err != nil {
			logger.Warn("failed to extract token from form", zap.Error(err))
			writeJSONError(w, http.StatusBadRequest, "missing or malformed access token", logger)
			return
		}

		jwksURL, err := GetJWKSURI(p.MetadataURL(), c)
		if err != nil {
			logger.Error("failed to fetch JWKS URI", zap.String("metadata_url", p.MetadataURL()), zap.Error(err))
			writeJSONError(w, http.StatusInternalServerError, "internal server error", logger)
			return
		}

		jwkSet, err := FetchJWKSet(jwksURL, c)
		if err != nil {
			logger.Error("failed to fetch JWKS", zap.String("metadata_url", jwksURL), zap.Error(err))
			writeJSONError(w, http.StatusInternalServerError, "internal server error", logger)
			return
		}

		kid, err := ExtractKIDFromToken(tokenStr)
		if err != nil {
			logger.Warn("failed to extract kid from token", zap.Error(err))
			writeJSONError(w, http.StatusBadRequest, "invalid token format", logger)
			return
		}

		jwk, err := FindJWKByKID(kid, jwkSet)
		if err != nil {
			logger.Warn("jwk not found for given kid", zap.Error(err))
			writeJSONError(w, http.StatusUnauthorized, "unauthorized", logger)
			return
		}

		pubKey, err := JWKToRSAPublicKey(jwk)
		if err != nil {
			logger.Error("failed to build RSA public key", zap.Error(err))
			writeJSONError(w, http.StatusInternalServerError, "internal server error", logger)
			return
		}

		claims, err := ParseAndVerifyJWT(tokenStr, pubKey, p.Issuer(), p.Audience())
		if err != nil {
			logger.Warn("JWT signature or claims validation failed", zap.Error(err))
			writeJSONError(w, http.StatusUnauthorized, "invalid token", logger)
			return
		}

		resp := UserResponse{
			Sub: claims.Subject,
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			logger.Error("failed to write user response", zap.Error(err))
		}
	}
}
