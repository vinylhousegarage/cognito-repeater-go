package deps

type RouteDependencies struct {
	CallbackProvider CallbackHandlerProvider
	LoginProvider    LoginHandlerProvider
	LogoutProvider   LogoutHandlerProvider
	MeProvider       MeHandlerProvider
	RevokeProvider   RevokeHandlerProvider
	WhoamiProvider   WhoamiHandlerProvider
}

type CallbackHandlerProvider interface {
	ClientSecretValue() string
	MetadataURL() string
	RedirectURIValue() string
	UserPoolClientIDValue() string
}

type LoginHandlerProvider interface {
	MetadataURL() string
	RedirectURIValue() string
	ScopeValue() string
	UserPoolClientIDValue() string
}

type LogoutHandlerProvider interface {
	MetadataURL() string
}

type MeHandlerProvider interface {
	Audience() string
	GetJWKSURI() string
	Issuer() string
	MetadataURL() string
}

type RevokeHandlerProvider interface {
	ClientSecretValue() string
	MetadataURL() string
	UserPoolClientIDValue() string
}

type WhoamiHandlerProvider interface {
	MetadataURL() string
}
