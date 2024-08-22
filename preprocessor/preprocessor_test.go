package preprocessor

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/pennsieve/processor-pre-external-files/models"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

func TestRun(t *testing.T) {
	integrationID := uuid.NewString()
	inputDir := t.TempDir()
	outputDir := t.TempDir()

	// For basic auth file download
	externalUsername := uuid.NewString()
	externalPassword := uuid.NewString()

	// For bearer auth file download
	externalToken := uuid.NewString()

	mock := NewMockServer()
	defer mock.Close()

	mockURL := mock.URL()

	externalFileParams := []models.ExternalFileParam{
		{
			URL:  fmt.Sprintf("%s/file1", mockURL),
			Name: "file1.png",
		},
		{
			URL:  fmt.Sprintf("%s/file2", mockURL),
			Name: "file2.json",
			Auth: &models.Auth{Type: models.BasicAuthType,
				Params: json.RawMessage(fmt.Sprintf(`{"username": %q, "password": %q}`, externalUsername, externalPassword)),
			},
		},
		{
			URL:  fmt.Sprintf("%s/file3", mockURL),
			Name: "file3.csv",
			Auth: &models.Auth{Type: models.BearerType,
				Params: json.RawMessage(fmt.Sprintf(`{"token": %q}`, externalToken)),
			},
			Query: map[string]string{"limit": "1000", "offset": "0"},
		},
	}
	// Create config file where pre-processor will expect it
	configFile, err := os.Create(filepath.Join(inputDir, ConfigFilename))
	require.NoError(t, err)
	require.NoError(t, json.NewEncoder(configFile).Encode(externalFileParams))

	expectedFiles := NewExpectedFiles(externalFileParams).Build(t, mockURL)
	mock.SetExpectedHandlers(t, expectedFiles)
	mock.Start()

	metadataPP := NewExternalFilesPreProcessor(integrationID, inputDir, outputDir)

	require.NoError(t, metadataPP.Run())
	expectedFiles.AssertEqual(t, inputDir)

}
