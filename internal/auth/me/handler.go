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

func NewMeHandler(
	p deps.MeHandlerProvider,
	c httpclient.HTTPClient,
	logger *zap.Logger,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Vary", "Origin")
			w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Authorization")
			w.WriteHeader(http.StatusNoContent)
			return
		}

		idToken, err := utils.ExtractAuthHeaderToken(r)
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
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Vary", "Origin")
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			logger.Error("failed to write user response", zap.Error(err))
			response.WriteErrorResponse(w, ErrFailedToWriteUserResponse, logger)
			return
		}

		logger.Info("user token verified successfully", zap.String("sub", claims.Subject))
	}
}
