package authtesthelpers

type MockMetadataURL struct {
	URL string
}

func (m *MockMetadataURL) MetadataURL() string {
	return m.URL
}

type MockMetadataURLProvider struct{}

func (m *MockMetadataURLProvider) MetadataURL() string {
	return "https://mock.metadata.url"
}
