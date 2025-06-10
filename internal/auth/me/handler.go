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
// @Failure 502 {object} response.ErrorResponse "Internal Bad Gateway"
// @Router /me [post]
func NewMeHandler(
	p deps.MeHandlerProvider,
	c httpclient.HTTPClient,
	logger *zap.Logger,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idToken, err := utils.ExtractFormValue(r)
		if err != nil {
			var status int
			switch {
			case errors.Is(err, utils.ErrFailedToParseForm),
				errors.Is(err, utils.ErrMissingToken):
				status = http.StatusBadRequest
				logger.Warn("ExtractFormValue returned an upstream error", zap.Error(err))
			default:
				status = http.StatusInternalServerError
				logger.Error("ExtractFormValue failed due to internal error", zap.Error(err))
			}
			utils.WritePlainError(w, status, err, logger)
			return
		}

		metadataURL := p.MetadataURL()
		endpoint, err := GetJWKSURI(metadataURL, c, logger)
		if err != nil {
			var status int
			switch {
			case errors.Is(err, ErrUnexpectedStatusCode),
				errors.Is(err, ErrMissingAuthorizationEndpoint):
				status = http.StatusBadGateway
				logger.Warn("GetJWKSURI returned an upstream error", zap.Error(err))
			default:
				status = http.StatusInternalServerError
				logger.Error("GetJWKSURI failed due to internal error", zap.Error(err))
			}
			utils.WritePlainError(w, status, err, logger)
			return
		}

		jwkSet, err := FetchJWKSet(jwksURL, c, logger)
		if err != nil {
			var status int
			var logMsg string
			switch {
			case errors.Is(err, ErrFailedToFetchJWKS):
				status = http.StatusBadGateway
				logMsg = "FetchJWKSet returned an upstream error"
			default:
				status = http.StatusInternalServerError
				logMsg = "FetchJWKSet failed due to internal error"
			}
			logger.Error(logMsg, zap.Error(err))
			utils.WritePlainError(w, status, err, logger)
			return
		}

		kid, err := ExtractKIDFromToken(idToken, logger)
		if err != nil {
			var status int
			switch {
			case errors.Is(err, ErrInvalidJWTFormat),
				errors.Is(err, ErrFailedToDecodeJWTHeader),
				errors.Is(err, ErrFailedToParseJWTHeader),
				errors.Is(err, ErrMissingKID):
				status = http.StatusBadRequest
				logger.Warn("ExtractKIDFromToken returned a client-side token error", zap.Error(err))
			default:
				status = http.StatusInternalServerError
				logger.Error("ExtractKIDFromToken failed due to an unexpected internal error", zap.Error(err))
			}
			utils.WritePlainError(w, status, err, logger)
			return
		}

		jwk, err := FindJWKByKID(kid, jwkSet, logger)
		if err != nil {
			var status int
			switch {
			case errors.Is(err, ErrJWKSetNil),
				errors.Is(err, ErrJWKNotFound):
				status = http.StatusUnauthorized
				logger.Warn("FindJWKByKID returned an authentication error", zap.Error(err))
			default:
				status = http.StatusInternalServerError
				logger.Error("FindJWKByKID failed due to an unexpected internal error", zap.Error(err))
			}
			utils.WritePlainError(w, status, err, logger)
			return
		}

		pubKey, err := JWKToRSAPublicKey(jwk, logger)
		if err != nil {
			var status int
			var logMsg string
			switch {
			case errors.Is(err, ErrInvalidN),
				errors.Is(err, ErrInvalidE):
				status = http.StatusInternalServerError
				logMsg = "JWKToRSAPublicKey failed due to invalid JWK components or internal error"
			default:
				status = http.StatusInternalServerError
				logMsg = "JWKToRSAPublicKey failed due to an unexpected internal error"
			}
			logger.Error(logMsg, zap.Error(err))
			utils.WritePlainError(w, status, err, logger)
			return
		}

		claims, err := ParseAndVerifyJWT(idToken, pubKey, p.Issuer(), p.Audience(), time.Now())
		if err != nil {
			var status int
			var logMsg string
			switch {
			case errors.Is(err, ErrJWTParseFailed),
				errors.Is(err, ErrInvalidSigningAlg),
				errors.Is(err, ErrTokenExpired),
				errors.Is(err, ErrInvalidIssuer),
				errors.Is(err, ErrMissingAudience),
				errors.Is(err, ErrInvalidAudience),
				errors.Is(err, ErrMissingSubject):
				status = http.StatusUnauthorized
				logMsg = "JWT validation failed due to authentication error"
			default:
				status = http.StatusInternalServerError
				logMsg = "JWT validation failed due to an unexpected internal error"
			}
			logger.Warn(logMsg, zap.Error(err))
			utils.WritePlainError(w, status, err, logger)
			return
		}

		resp := UserResponse{
			Sub: claims.Subject,
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			logger.Error("failed to write user response", zap.Error(err))
			http.Error(w, "failed to write response", http.StatusInternalServerError)
		}

		logger.Info("user token verified successfully", zap.String("sub", claims.Subject))
	}
}
