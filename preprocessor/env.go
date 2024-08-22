package preprocessor

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
	return NewExternalFilesPreProcessor(integrationID, inputDirectory, outputDirectory), nil
}
