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
	return NewExternalFilesPreProcessor(integrationID, inputDirectory, outputDirectory, externalFiles), nil
}
