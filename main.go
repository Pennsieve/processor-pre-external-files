package main

import (
	"github.com/pennsieve/processor-pre-external-files/logging"
	"github.com/pennsieve/processor-pre-external-files/preprocessor"
	"log/slog"
	"os"
)

var logger = logging.PackageLogger("main")

func main() {

	m, err := preprocessor.FromEnv()
	if err != nil {
		logger.Error("error creating preprocessor", slog.Any("error", err))
		os.Exit(1)
	}

	logger.Info("created ExternalFilesPreProcessor",
		slog.String("integrationID", m.IntegrationID),
		slog.String("inputDirectory", m.InputDirectory),
		slog.String("outputDirectory", m.OutputDirectory),
		slog.Any("externalFiles", m.ExternalFiles),
	)

	if err := m.Run(); err != nil {
		logger.Error("error running preprocessor", slog.Any("error", err))
		os.Exit(1)
	}
}
