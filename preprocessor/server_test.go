package preprocessor

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/pennsieve/processor-pre-external-files/models"
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

// URL is necessary because we need the URL before adding handlers, so before the server starts
// and httptest.Server.URL is empty before it starts. This work-around taken from httptest.NewServer() code.
func (m *MockServer) URL() string {
	return fmt.Sprintf("http://%s", m.Server.Listener.Addr().String())
}

func (m *MockServer) SetExpectedHandlers(t *testing.T, integrationID string, expectedFiles *ExpectedFiles) {
	m.Mux.HandleFunc(fmt.Sprintf("/integrations/%s", integrationID), func(writer http.ResponseWriter, request *http.Request) {
		require.Equal(t, http.MethodGet, request.Method, "expected method %s for %s, got %s", http.MethodGet, request.URL, request.Method)
		integration := models.Integration{
			Uuid:          uuid.NewString(),
			ApplicationID: 0,
			PackageIDs:    nil,
			Params:        models.IntegrationParams{ExternalFiles: expectedFiles.ExternalFileParams},
		}
		integrationResponse, err := json.Marshal(integration)
		require.NoError(t, err)
		_, err = writer.Write(integrationResponse)
		require.NoError(t, err)
	})
	for _, expectedFile := range expectedFiles.Files {
		m.Mux.HandleFunc(expectedFile.APIPath, expectedFile.HandlerFunc(t))
	}
	m.Mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		require.Fail(t, "unexpected call to Pennsieve", "%s %s", request.Method, request.URL)
	})
}
