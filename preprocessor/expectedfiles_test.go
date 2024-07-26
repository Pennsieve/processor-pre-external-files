package preprocessor

import (
	crypto "crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/pennsieve/processor-pre-external-files/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// No tests, just test helpers

type ExpectedFile struct {
	Name  string
	Bytes []byte
	// APIPath is the request path the mock server will match against.
	APIPath             string
	QueryParams         url.Values
	AuthorizationHeader string
}

func FromExternalFileParam(t *testing.T, externalFileParam models.ExternalFileParam, urlPrefix string) ExpectedFile {
	toStrip := urlPrefix
	if strings.HasSuffix(toStrip, "/") {
		toStrip = strings.TrimSuffix(toStrip, "/")
	}
	require.True(t, strings.HasPrefix(externalFileParam.URL, toStrip), "URL %s does not start with expected prefix %s", externalFileParam.URL, toStrip)

	apiPath := strings.TrimPrefix(externalFileParam.URL, toStrip)

	size := rand.Intn(1000) + 1
	bytes := make([]byte, size)
	_, err := crypto.Read(bytes)
	require.NoError(t, err)

	urlValues := url.Values{}
	for k, v := range externalFileParam.Query {
		urlValues.Set(k, v)
	}

	authHeader := ""
	if externalFileParam.Auth != nil {
		switch externalFileParam.Auth.Type {
		case models.BasicAuthType:
			var params models.BasicAuthParams
			require.NoError(t, externalFileParam.Auth.UnmarshalParams(&params))
			toEncode := fmt.Sprintf("%s:%s", params.Username, params.Password)
			authHeader = fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(toEncode)))
		case models.BearerType:
			var params models.BearerAuthParams
			require.NoError(t, externalFileParam.Auth.UnmarshalParams(&params))
			authHeader = fmt.Sprintf("Bearer %s", params.Token)
		default:
			require.Fail(t, "unknown auth type", externalFileParam.Auth.Type)
		}

	}

	return ExpectedFile{
		Name:                externalFileParam.Name,
		Bytes:               bytes,
		APIPath:             apiPath,
		QueryParams:         urlValues,
		AuthorizationHeader: authHeader,
	}
}

func (e ExpectedFile) HandlerFunc(t *testing.T) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		require.Equal(t, http.MethodGet, request.Method, "expected method %s for %s, got %s", http.MethodGet, request.URL, request.Method)
		actualAuthorizationHeader := request.Header.Get("authorization")
		require.Equal(t, e.AuthorizationHeader, actualAuthorizationHeader)

		if len(e.QueryParams) > 0 {
			require.Equal(t, e.QueryParams, request.URL.Query(), "expected query %s for %s, got %s", e.QueryParams, request.URL, request.URL.Query())
		}
		_, err := writer.Write(e.Bytes)
		require.NoError(t, err)
	}
}

type ExpectedFiles struct {
	ExternalFileParams []models.ExternalFileParam
	DatasetID          string
	Files              []ExpectedFile
}

func NewExpectedFiles(externalFileParams []models.ExternalFileParam) *ExpectedFiles {
	return &ExpectedFiles{
		ExternalFileParams: externalFileParams,
	}
}

func (e *ExpectedFiles) Build(t *testing.T, mockURL string) *ExpectedFiles {
	for _, external := range e.ExternalFileParams {
		e.Files = append(e.Files, FromExternalFileParam(t, external, mockURL))
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
