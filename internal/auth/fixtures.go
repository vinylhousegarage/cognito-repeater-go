package auth

type mockMetadataURL struct {
	URL string
}

func (m *mockMetadataURL) MetadataURL() string {
	return m.URL
}
