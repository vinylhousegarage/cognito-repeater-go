package authtesthelpers

type MockMetadataURL struct {
	URL string
}

func (m *MockMetadataURL) MetadataURL() string {
	return m.URL
}

type MmockEndpointProvider struct{}

func (m *MockEndpointProvider) GetLogoutURL(p config.MetadataURLProvider) (string, error) {
	return "https://example.com/logout", nil
}

type MockMetadataProvider struct{}

func (m *MockMetadataProvider) MetadataURL() string {
	return "https://mock.metadata.url"
}
