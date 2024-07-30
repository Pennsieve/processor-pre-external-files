package preprocessor

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

// No tests, just test helpers

type MockServer struct {
	Mux           *http.ServeMux
	Server        *httptest.Server
	ExpectedFiles *ExpectedFiles
}

func NewMockServer() *MockServer {
	mock := &MockServer{}
	mock.Mux = http.NewServeMux()
	mock.Server = httptest.NewUnstartedServer(mock.Mux)
	return mock
}

func (m *MockServer) Start() {
	m.Server.Start()
}

func (m *MockServer) Close() {
	m.Server.Close()
}

// URL() function is necessary because we need the server's URL before adding handlers, that is, before the server starts.
// But httptest.Server.URL is empty before server starts. This work-around taken from httptest.NewServer() code.
func (m *MockServer) URL() string {
	return fmt.Sprintf("http://%s", m.Server.Listener.Addr().String())
}

func (m *MockServer) SetExpectedHandlers(t *testing.T, expectedFiles *ExpectedFiles) {
	for _, expectedFile := range expectedFiles.Files {
		m.Mux.HandleFunc(expectedFile.APIPath, expectedFile.HandlerFunc(t))
	}
	m.Mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		require.Fail(t, "unexpected call to Pennsieve", "%s %s", request.Method, request.URL)
	})
}
