package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

func initializeLogger(logFile string) (*log.Logger, error) {

	if logFile != "" {
		file, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o644)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file: %w", err)
		}
		multiWriter := io.MultiWriter(os.Stderr, file)
		return log.New(multiWriter, "", log.LstdFlags), nil
	}
	return log.New(os.Stderr, "", log.LstdFlags), nil
}
