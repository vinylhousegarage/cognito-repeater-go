package authtesthelpers

type MockMetadataURL struct {
	URL string
}

func (m *MockMetadataURL) MetadataURL() string {
	return m.URL
}

type MockMetadataURLProvider struct{}

func (p *MockMetadataURLProvider) MetadataURL() string {
	return "https://mock.metadata.url"
}
