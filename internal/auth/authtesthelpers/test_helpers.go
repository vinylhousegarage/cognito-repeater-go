package authtesthelpers

import (
	"net/http"
	"io"
	"strings"
)

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

type MockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

var MockClient = &MockHTTPClient{
	DoFunc: func(req *http.Request) (*http.Response, error) {
		body := `{"end_session_endpoint":"https://example.com/logout"}`
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(body)),
			Header:     make(http.Header),
		}, nil
	},
}
