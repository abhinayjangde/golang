package logging

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
)

func New(filePath string) (*slog.Logger, io.Closer, error) {
	var output io.Writer = os.Stdout
	var file *os.File

	if filePath != "" {
		if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
			return nil, nil, err
		}

		var err error
		file, err = os.OpenFile(
			filePath,
			os.O_CREATE|os.O_APPEND|os.O_WRONLY,
			0644,
		)
		if err != nil {
			return nil, nil, err
		}

		output = io.MultiWriter(os.Stdout, file)
	}

	handler := slog.NewJSONHandler(output, &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true,
	})

	logger := slog.New(handler)

	if file == nil {
		return logger, io.NopCloser(nil), nil
	}

	return logger, file, nil
}
