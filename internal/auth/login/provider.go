type LoginHandlerProvider interface {
	MetadataURL() string
	RedirectURIValue() string
	ScopeValue() string
	UserPoolClientIDValue() string
}
