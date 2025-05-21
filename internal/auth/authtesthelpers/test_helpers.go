package authtesthelpers

type MockMetadataURL struct {
	URL string
}

func (m *MockMetadataURL) MetadataURL() string {
	return m.URL
}

type MockMetadataProvider struct{}

func (m *MockMetadataProvider) MetadataURL() string {
	return "https://mock.metadata.url"
}
