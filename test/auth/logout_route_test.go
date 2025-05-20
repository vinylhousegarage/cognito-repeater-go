func TestLogoutRoute_IsRegistered(t *testing.T) {
	router := http.NewServeMux()
	router.Handle("/logout/redirect", LogoutHandler(&mockEndpointProvider{})(&mockMetadataProvider{}))

	req := httptest.NewRequest(http.MethodGet, "/logout/redirect", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusFound, resp.StatusCode)
}
