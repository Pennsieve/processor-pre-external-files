package preprocessor

import (
	"encoding/json"
	"fmt"
	"github.com/pennsieve/processor-pre-external-files/client/models"
	"github.com/pennsieve/processor-pre-external-files/service/logging"
	"github.com/pennsieve/processor-pre-external-files/service/util"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

var logger = logging.PackageLogger("preprocessor")

type ExternalFilesPreProcessor struct {
	IntegrationID   string
	InputDirectory  string
	OutputDirectory string
	ConfigFile      string
}

func NewExternalFilesPreProcessor(integrationID string,
	inputDirectory string,
	outputDirectory string,
	configFile string) *ExternalFilesPreProcessor {
	return &ExternalFilesPreProcessor{
		IntegrationID:   integrationID,
		InputDirectory:  inputDirectory,
		OutputDirectory: outputDirectory,
		ConfigFile:      configFile,
	}
}

func (m *ExternalFilesPreProcessor) Run() error {
	logger.Info("processing integration", slog.String("integrationID", m.IntegrationID))
	externalFiles, err := readConfig(m.ConfigFile)
	if err != nil {
		return err
	}

	if len(externalFiles) == 0 {
		logger.Info("integration contained no external files")
		return nil
	}
	for _, externalFile := range externalFiles {
		efLogger := externalFile.Logger(logger)
		efLogger.Info("handling external file")
		request, err := newRequest(externalFile)
		if err != nil {
			return err
		}
		response, err := util.Invoke(request)
		if err != nil {
			return err
		}
		downloadPath := filepath.Join(m.InputDirectory, externalFile.Name)
		written, err := writeResponse(response, downloadPath)
		if err != nil {
			return err
		}
		efLogger.Info("wrote file",
			slog.String("path", downloadPath),
			slog.Int64("size", written))
	}
	logger.Info("downloads complete")

	return nil
}

func readConfig(configPath string) (models.ExternalFileParams, error) {
	var config models.ExternalFileParams
	configFile, err := os.Open(configPath)
	if err != nil {
		return config, fmt.Errorf("error opening config file %s: %w", configPath, err)
	}
	if err := json.NewDecoder(configFile).Decode(&config); err != nil {
		return config, fmt.Errorf("error decoding config file %s: %w", configPath, err)
	}
	return config, nil
}

func newRequest(externalFile models.ExternalFileParam) (*http.Request, error) {
	fullURL, err := url.Parse(externalFile.URL)
	if err != nil {
		return nil, fmt.Errorf("error parsing URL %s: %w", externalFile.URL, err)
	}
	urlQuery := url.Values{}
	for q, v := range externalFile.Query {
		urlQuery.Add(q, v)
	}
	fullURL.RawQuery = urlQuery.Encode()
	req, err := http.NewRequest(http.MethodGet, fullURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request for %s: %w", externalFile.URL, err)
	}
	if err := externalFile.Auth.SetAuthentication(req); err != nil {
		return nil, err
	}
	return req, nil
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
