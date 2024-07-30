package preprocessor

import (
	"encoding/json"
	"fmt"
	"github.com/pennsieve/processor-pre-external-files/models"
)

const IntegrationIDKey = "INTEGRATION_ID"
const InputDirectoryKey = "INPUT_DIR"
const OutputDirectoryKey = "OUTPUT_DIR"
const ExternalFilesKey = "EXTERNAL_FILES"
const SessionTokenKey = "SESSION_TOKEN"
const PennsieveAPIHostKey = "PENNSIEVE_API_HOST"
const PennsieveAPI2HostKey = "PENNSIEVE_API_HOST2"

func FromEnv() (*ExternalFilesPreProcessor, error) {
	integrationID, err := LookupRequiredEnvVar(IntegrationIDKey)
	if err != nil {
		return nil, err
	}
	inputDirectory, err := LookupRequiredEnvVar(InputDirectoryKey)
	if err != nil {
		return nil, err
	}
	outputDirectory, err := LookupRequiredEnvVar(OutputDirectoryKey)
	if err != nil {
		return nil, err
	}
	externalFileParams, err := LookupRequiredEnvVar(ExternalFilesKey)
	if err != nil {
		return nil, err
	}
	var externalFiles models.ExternalFileParams
	if err := json.Unmarshal([]byte(externalFileParams), &externalFiles); err != nil {
		return nil, fmt.Errorf("error unmarshalling %s value %q: %w", ExternalFilesKey, externalFileParams, err)
	}
	sessionToken, err := LookupRequiredEnvVar(SessionTokenKey)
	if err != nil {
		return nil, err
	}
	apiHost, err := LookupRequiredEnvVar(PennsieveAPIHostKey)
	if err != nil {
		return nil, err
	}
	api2Host, err := LookupRequiredEnvVar(PennsieveAPI2HostKey)
	if err != nil {
		return nil, err
	}
	return NewExternalFilesPreProcessor(integrationID, inputDirectory, outputDirectory, externalFiles, sessionToken, apiHost, api2Host), nil
}
