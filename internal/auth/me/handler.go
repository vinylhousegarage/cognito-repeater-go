package me

import (
	"encoding/json"
	"log"
	"net/http"

	"cognito-repeater-go/internal/auth/utils"
	"cognito-repeater-go/internal/httpclient"
)

type MeHandlerProvider interface {
	Audience() string
	GetJWKSURI() string
	Issuer() string
	MetadataURL() string
}

func writeJSONError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"error": msg,
	})
}

func MeHandler(p MeHandlerProvider, c httpclient.HTTPClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenStr, err := utils.ExtractAuthHeaderToken(r)
		if err != nil {
			log.Printf("token extraction error: %v", err)
			writeJSONError(w, http.StatusBadRequest, "missing or malformed access token")
			return
		}

		jwksURL, err := GetJWKSURI(p.MetadataURL(), c)
		if err != nil {
			log.Printf("failed to fetch JWKS URI: %v", err)
			writeJSONError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		jwkSet, err := FetchJWKSet(jwksURL, c)
		if err != nil {
			log.Printf("failed to fetch JWKS: %v", err)
			writeJSONError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		kid, err := ExtractKIDFromToken(tokenStr)
		if err != nil {
			log.Printf("failed to extract kid: %v", err)
			writeJSONError(w, http.StatusBadRequest, "invalid token format")
			return
		}

		jwk, err := FindJWKByKID(kid, jwkSet)
		if err != nil {
			log.Printf("jwk not found: %v", err)
			writeJSONError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		pubKey, err := JWKToRSAPublicKey(jwk)
		if err != nil {
			log.Printf("failed to build RSA public key: %v", err)
			writeJSONError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		claims, err := ParseAndVerifyJWT(tokenStr, pubKey, p.Issuer(), p.Audience())
		if err != nil {
			log.Printf("token verification failed: %v", err)
			writeJSONError(w, http.StatusUnauthorized, "invalid token")
			return
		}

		if err := json.NewEncoder(w).Encode(map[string]string{
			"sub": claims.Subject,
		}); err != nil {
			log.Printf("failed to write response: %v", err)
		}
	}
}
