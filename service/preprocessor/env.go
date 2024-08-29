package preprocessor

import "path/filepath"

const IntegrationIDKey = "INTEGRATION_ID"
const InputDirectoryKey = "INPUT_DIR"
const OutputDirectoryKey = "OUTPUT_DIR"
const ConfigDirectoryKey = "CONFIG_DIR"

const DefaultConfigFilename = "config.json"

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
	configDirectory, err := LookupRequiredEnvVar(ConfigDirectoryKey)
	if err != nil {
		return nil, err
	}
	configFile := filepath.Join(configDirectory, DefaultConfigFilename)
	return NewExternalFilesPreProcessor(integrationID, inputDirectory, outputDirectory, configFile), nil
}
