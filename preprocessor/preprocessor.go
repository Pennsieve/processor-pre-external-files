package preprocessor

import (
	"fmt"
	"github.com/pennsieve/processor-pre-external-files/logging"
	"github.com/pennsieve/processor-pre-external-files/pennsieve"
	"github.com/pennsieve/processor-pre-external-files/util"
	"io"
	"log/slog"
	"net/http"
	"os"
)

var logger = logging.PackageLogger("preprocessor")

type ExternalFilesPreProcessor struct {
	IntegrationID   string
	InputDirectory  string
	OutputDirectory string
	Pennsieve       *pennsieve.Session
}

func NewExternalFilesPreProcessor(integrationID string,
	inputDirectory string,
	outputDirectory string,
	sessionToken string,
	apiHost string,
	api2Host string) *ExternalFilesPreProcessor {
	return &ExternalFilesPreProcessor{
		IntegrationID:   integrationID,
		InputDirectory:  inputDirectory,
		OutputDirectory: outputDirectory,
		Pennsieve:       pennsieve.NewSession(sessionToken, apiHost, api2Host),
	}
}

func FromEnv() (*ExternalFilesPreProcessor, error) {
	integrationID, err := LookupRequiredEnvVar("INTEGRATION_ID")
	if err != nil {
		return nil, err
	}
	inputDirectory, err := LookupRequiredEnvVar("INPUT_DIR")
	if err != nil {
		return nil, err
	}
	outputDirectory, err := LookupRequiredEnvVar("OUTPUT_DIR")
	if err != nil {
		return nil, err
	}
	sessionToken, err := LookupRequiredEnvVar("SESSION_TOKEN")
	if err != nil {
		return nil, err
	}
	apiHost, err := LookupRequiredEnvVar("PENNSIEVE_API_HOST")
	if err != nil {
		return nil, err
	}
	api2Host, err := LookupRequiredEnvVar("PENNSIEVE_API_HOST2")
	if err != nil {
		return nil, err
	}
	return NewExternalFilesPreProcessor(integrationID, inputDirectory, outputDirectory, sessionToken, apiHost, api2Host), nil
}

func (m *ExternalFilesPreProcessor) Run() error {
	// get integration info
	integration, err := m.Pennsieve.GetIntegration(m.IntegrationID)
	if err != nil {
		return err
	}
	logger.Info("got integration", slog.Any("params", integration.Params),
		slog.Any("workflow", integration.Workflow))

	return nil
}

func writeResponse(response *http.Response, filePath string) (int64, error) {
	defer util.CloseAndWarn(response)

	file, err := os.Create(filePath)
	if err != nil {
		return 0, fmt.Errorf("error creating file %s: %w", filePath, err)
	}
	written, err := io.Copy(file, response.Body)
	if err != nil {
		return 0, fmt.Errorf("error writing %s %s response to %s: %w",
			response.Request.Method,
			response.Request.URL,
			filePath,
			err)
	}
	return written, nil
}

func LookupRequiredEnvVar(key string) (string, error) {
	value := os.Getenv(key)
	if len(value) == 0 {
		return "", fmt.Errorf("no %s set", key)
	}
	return value, nil
}
