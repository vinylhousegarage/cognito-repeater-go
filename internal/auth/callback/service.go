package callback

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"cognito-repeater-go/internal/config"
	"cognito-repeater-go/internal/httpclient"
)

func ValidateCallbackRequest(r *http.Request) (string, error) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	if code == "" {
		return "", errors.New("missing code")
	}

	if state == "" {
		return "", errors.New("missing state")
	}

	cookie, err := r.Cookie("oauth_state")
	if err != nil {
		return "", errors.New("missing oauth_state cookie")
	}

	if state != cookie.Value {
		return "", errors.New("invalid state")
	}

	return code, nil
}

type callbackMetadata struct {
	TokenEndpoint string `json:"token_endpoint"`
}

type CallbackURLProvider interface {
	GetCallbackURL(p config.MetadataURLProvider) (string, error)
}

type callbackClient struct {
	client httpclient.HTTPClient
}

func NewCallbackClient(client httpclient.HTTPClient) CallbackURLProvider {
	return &callbackClient{client: client}
}

func (s *callbackClient) GetCallbackURL(p config.MetadataURLProvider) (string, error) {
	url := p.MetadataURL()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch metadata: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var meta callbackMetadata
	if err := json.NewDecoder(resp.Body).Decode(&meta); err != nil {
		return "", fmt.Errorf("failed to decode metadata: %w", err)
	}

	if meta.TokenEndpoint == "" {
		return "", fmt.Errorf("invalid metadata: token endpoint is empty")
	}

	return meta.TokenEndpoint, nil
}
