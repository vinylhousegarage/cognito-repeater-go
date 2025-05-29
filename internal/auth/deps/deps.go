package deps

import (
	"cognito-repeater-go/internal/config"
	"cognito-repeater-go/internal/httpclient"
)

type HandlerDependencies struct {
	Config     config.CognitoMetadataProvider
	HTTPClient httpclient.HTTPClient
}
