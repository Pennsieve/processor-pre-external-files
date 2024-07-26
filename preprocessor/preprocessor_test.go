package preprocessor

import (
	crypto "crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/pennsieve/processor-pre-external-files/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"testing"
)

func TestRun(t *testing.T) {
	datasetId := uuid.NewString()

	integrationID := uuid.NewString()
	inputDir := t.TempDir()
	outputDir := t.TempDir()
	sessionToken := uuid.NewString()
	username := uuid.NewString()
	password := uuid.NewString()

	externalFileParams := []models.ExternalFileParam{
		{
			URL: "https://example.com/file1",
		},
		{
			URL: "https://example.com/file2",
			Auth: &models.Auth{Type: models.BasicAuthType,
				Params: json.RawMessage(fmt.Sprintf(`{"username": %q, "password": %q}`, username, password)),
			},
		},
	}
	expectedFiles := NewExpectedFiles(datasetId).Build(t)

	mockServer := newMockServer(t, integrationID, externalFileParams, expectedFiles)
	defer mockServer.Close()

	mockServer.Start()

	metadataPP := NewExternalFilesPreProcessor(integrationID, inputDir, outputDir, sessionToken, mockServer.URL, mockServer.URL)

	require.NoError(t, metadataPP.Run())
	expectedFiles.AssertEqual(t, inputDir)

}

type ExpectedFile struct {
	Name  string
	Bytes []byte
	// APIPath is the request path the mock server will match against.
	APIPath     string
	QueryParams url.Values
}

func (e ExpectedFile) HandlerFunc(t *testing.T) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		require.Equal(t, http.MethodGet, request.Method, "expected method %s for %s, got %s", http.MethodGet, request.URL, request.Method)
		if e.QueryParams != nil {
			require.Equal(t, e.QueryParams, request.URL.Query(), "expected query %s for %s, got %s", e.QueryParams, request.URL, request.URL.Query())
		}
		_, err := writer.Write(e.Bytes)
		require.NoError(t, err)
	}
}

type ExpectedFiles struct {
	DatasetID string
	Files     []ExpectedFile
}

func NewExpectedFiles(datasetID string) *ExpectedFiles {
	return &ExpectedFiles{
		DatasetID: datasetID,
	}
}

func (e *ExpectedFiles) Build(t *testing.T) *ExpectedFiles {
	for i := range e.Files {
		expected := &e.Files[i]
		size := rand.Intn(1000) + 1
		bytes := make([]byte, size)
		_, err := crypto.Read(bytes)
		require.NoError(t, err)
		expected.Bytes = bytes
	}
	return e
}

func (e *ExpectedFiles) AssertEqual(t *testing.T, actualDir string) {
	for _, expectedFile := range e.Files {
		actualFilePath := filepath.Join(actualDir, expectedFile.Name)
		actualBytes, err := os.ReadFile(actualFilePath)
		if assert.NoError(t, err) {
			assert.Equal(t, expectedFile.Bytes, actualBytes, "actual bytes %s do not match expected bytes %s", actualFilePath, expectedFile.Name)
		}
	}
}

func newMockServer(t *testing.T, integrationID string, externalFilesParams []models.ExternalFileParam, expectedFiles *ExpectedFiles) *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc(fmt.Sprintf("/integrations/%s", integrationID), func(writer http.ResponseWriter, request *http.Request) {
		require.Equal(t, http.MethodGet, request.Method, "expected method %s for %s, got %s", http.MethodGet, request.URL, request.Method)
		integration := models.Integration{
			Uuid:          uuid.NewString(),
			ApplicationID: 0,
			PackageIDs:    nil,
			Params:        models.IntegrationParams{ExternalFiles: externalFilesParams},
		}
		integrationResponse, err := json.Marshal(integration)
		require.NoError(t, err)
		_, err = writer.Write(integrationResponse)
		require.NoError(t, err)
	})
	for _, expectedFile := range expectedFiles.Files {
		mux.HandleFunc(expectedFile.APIPath, expectedFile.HandlerFunc(t))
	}
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		require.Fail(t, "unexpected call to Pennsieve", "%s %s", request.Method, request.URL)
	})
	return httptest.NewUnstartedServer(mux)
}
