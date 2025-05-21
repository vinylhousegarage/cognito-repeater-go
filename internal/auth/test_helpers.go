package auth

type MockMetadataURL struct {
	URL string
}

func (m *MockMetadataURL) MetadataURL() string {
	return m.URL
}
