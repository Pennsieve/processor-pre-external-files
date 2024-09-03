package models

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestExternalFilesParamUnmarshal_Just_URL(t *testing.T) {
	expectedURL := "https://example.com/file"
	expectedName := "test.txt"
	jsonString := fmt.Sprintf(`{"url": %q, "name": %q}`, expectedURL, expectedName)
	var params ExternalFileParam
	require.NoError(t, json.Unmarshal([]byte(jsonString), &params))
	assert.Equal(t, expectedURL, params.URL)
	assert.Equal(t, expectedName, params.Name)
	assert.Nil(t, params.Auth)
	assert.Empty(t, params.Query)
}

func TestExternalFilesParamUnmarshal_Query(t *testing.T) {
	expectedURL := "https://example.com/file"
	expectedName := "test.txt"
	expectedKey := "datasetId"
	expectedValue := "N:dataset:123-456"
	expectedKey2 := "id"
	expectedValue2 := "6"
	jsonString := fmt.Sprintf(`{
									"url": %q,
									"name": %q,
									"query": {%q: %q, %q: %q}
								}`,
		expectedURL,
		expectedName,
		expectedKey,
		expectedValue,
		expectedKey2,
		expectedValue2)
	var params ExternalFileParam
	require.NoError(t, json.Unmarshal([]byte(jsonString), &params))
	assert.Equal(t, expectedURL, params.URL)
	assert.Equal(t, expectedName, params.Name)
	assert.Nil(t, params.Auth)

	query := params.Query
	assert.Len(t, query, 2)

	assert.Contains(t, query, expectedKey)
	assert.Equal(t, expectedValue, query[expectedKey])

	assert.Contains(t, query, expectedKey2)
	assert.Equal(t, expectedValue2, query[expectedKey2])
}

func TestExternalFilesParamUnmarshal_Basic_Auth(t *testing.T) {
	expectedURL := "https://example.com/file"
	expectedName := "test.txt"
	expectedUsername := "joe"
	expectedPassword := uuid.NewString()
	jsonString := fmt.Sprintf(`{
									"url": %q,
									"name": %q,
									"auth": {"type": "basicAuth", "params": {"username": %q, "password": %q}}
								}`,
		expectedURL,
		expectedName,
		expectedUsername,
		expectedPassword)
	var params ExternalFileParam
	require.NoError(t, json.Unmarshal([]byte(jsonString), &params))
	assert.Equal(t, expectedURL, params.URL)
	assert.Equal(t, expectedName, params.Name)
	assert.Empty(t, params.Query)

	assert.NotNil(t, params.Auth)

	request, err := http.NewRequest(http.MethodGet, expectedURL, nil)
	require.NoError(t, err)

	require.NoError(t, params.Auth.SetAuthentication(request))
	actualUsername, actualPassword, isBasicAuth := request.BasicAuth()
	if assert.True(t, isBasicAuth) {
		assert.Equal(t, expectedUsername, actualUsername)
		assert.Equal(t, expectedPassword, actualPassword)
	}

}

func TestExternalFilesParamUnmarshal_Bearer_Auth(t *testing.T) {
	expectedURL := "https://example.com/file"
	expectedName := "test.txt"
	expectedToken := uuid.NewString()
	jsonString := fmt.Sprintf(`{
									"url": %q,
									"name": %q,
									"auth": {"type": "bearer", "params": {"token": %q}}
								}`,
		expectedURL,
		expectedName,
		expectedToken)
	var params ExternalFileParam
	require.NoError(t, json.Unmarshal([]byte(jsonString), &params))
	assert.Equal(t, expectedURL, params.URL)
	assert.Equal(t, expectedName, params.Name)
	assert.Empty(t, params.Query)

	assert.NotNil(t, params.Auth)

	request, err := http.NewRequest(http.MethodGet, expectedURL, nil)
	require.NoError(t, err)

	require.NoError(t, params.Auth.SetAuthentication(request))

	actualAuth := request.Header.Get("authorization")
	assert.Equal(t, fmt.Sprintf("Bearer %s", expectedToken), actualAuth)

}
