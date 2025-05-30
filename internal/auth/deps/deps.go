package deps

type RouteDependencies struct {
	CallbackProvider callback.CallbackHandlerProvider
	LoginProvider    login.LoginHandlerProvider
	LogoutProvider   logout.LogoutHandlerProvider
	MeProvider       me.MeHandlerProvider
}
