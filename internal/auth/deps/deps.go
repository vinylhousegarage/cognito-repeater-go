package deps

import (
	"cognito-repeater-go/internal/auth/callback"
	"cognito-repeater-go/internal/auth/login"
	"cognito-repeater-go/internal/auth/logout"
	"cognito-repeater-go/internal/auth/me"
)

type RouteDependencies struct {
	CallbackProvider callback.CallbackHandlerProvider
	LoginProvider    login.LoginHandlerProvider
	LogoutProvider   logout.LogoutHandlerProvider
	MeProvider       me.MeHandlerProvider
}
