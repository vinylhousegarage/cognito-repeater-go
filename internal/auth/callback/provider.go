package callback

type CallbackHandlerProvider interface {
	ClientSecretValue() string
	MetadataURL() string
	RedirectURIValue() string
	UserPoolClientIDValue() string
}
