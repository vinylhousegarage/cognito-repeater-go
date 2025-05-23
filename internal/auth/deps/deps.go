package deps

import (
	"cognito-repeater-go/internal/auth/config"
	"cognito-repeater-go/internal/httpclient"
)

type HandlerDependencies struct {
	Config     *config.Config
	HTTPClient httpclient.HTTPClient
}
