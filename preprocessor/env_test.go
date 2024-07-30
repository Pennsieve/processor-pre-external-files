package preprocessor

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/pennsieve/processor-pre-external-files/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFromEnv(t *testing.T) {

	expectedExternalFiles := models.ExternalFileParams{
		{
			URL:   "https://httpbin.org/get",
			Name:  "no-auth-with-query.json",
			Query: map[string]string{"param1": "9", "param2": "xyz"},
		},
		{
			URL:  "https://httpbin.org/basic-auth/joe/JoE123",
			Name: "basic-auth.json",
			Auth: &models.Auth{Type: models.BasicAuthType, Params: json.RawMessage(`{"username": "joe", "password": "JoE123"}`)},
		},
		{
			URL:  "https://httpbin.org/bearer",
			Name: "bearer-auth.json",
			Auth: &models.Auth{Type: models.BearerType, Params: json.RawMessage(`{"token": "JoE123-api-key"}`)},
		},
	}
	externalFilesBytes, err := json.Marshal(expectedExternalFiles)
	require.NoError(t, err)
	fmt.Println(string(externalFilesBytes))
	t.Setenv(ExternalFilesKey, string(externalFilesBytes))

	expectedSessionToken := uuid.NewString()
	t.Setenv(SessionTokenKey, expectedSessionToken)
	expectedIntegrationID := uuid.NewString()
	t.Setenv(IntegrationIDKey, expectedIntegrationID)
	expectedInputDirectory := fmt.Sprintf("input/%s", uuid.NewString())
	t.Setenv(InputDirectoryKey, expectedInputDirectory)
	expectedOutputDirectory := fmt.Sprintf("output/%s", uuid.NewString())
	t.Setenv(OutputDirectoryKey, expectedOutputDirectory)
	expectedAPIHost := "https://pennsieve.example.com"
	t.Setenv(PennsieveAPIHostKey, expectedAPIHost)
	expectedAPI2Host := "https://pennsieve2.example.com"
	t.Setenv(PennsieveAPI2HostKey, expectedAPI2Host)

	processor, err := FromEnv()
	require.NoError(t, err)

	assert.Equal(t, expectedSessionToken, processor.Pennsieve.Token)
	assert.Equal(t, expectedIntegrationID, processor.IntegrationID)
	assert.Equal(t, expectedInputDirectory, processor.InputDirectory)
	assert.Equal(t, expectedOutputDirectory, processor.OutputDirectory)
	assert.Equal(t, expectedAPIHost, processor.Pennsieve.APIHost)
	assert.Equal(t, expectedAPI2Host, processor.Pennsieve.API2Host)

	actualExternalFiles := processor.ExternalFiles
	assert.NotNil(t, actualExternalFiles)
	assert.Len(t, actualExternalFiles, len(expectedExternalFiles))
	// Can't just assert.Equal because json.RawMessage fields might have been serialized with different whitespace
	for i := 0; i < len(actualExternalFiles); i++ {
		expected := expectedExternalFiles[i]
		actual := actualExternalFiles[i]
		assert.Equal(t, expected.URL, actual.URL)
		assert.Equal(t, expected.Query, actual.Query)
		assert.Equal(t, expected.Name, actual.Name)
		if expected.Auth == nil {
			assert.Nil(t, actual.Auth)
		} else {
			assert.Equal(t, expected.Auth.Type, actual.Auth.Type)
			expectedAuthParams, err := expected.Auth.UnmarshallParams()
			require.NoError(t, err)
			actualAuthParams, err := actual.Auth.UnmarshallParams()
			require.NoError(t, err)
			assert.Equal(t, expectedAuthParams, actualAuthParams)
		}
	}
}
