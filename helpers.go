package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

type closeFunc func() error

func initializeLogger(logFile string) (*log.Logger, closeFunc, error) {

	if logFile != "" {
		file, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o644)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to open log file: %w", err)
		}
		bufferedFile := bufio.NewWriterSize(file, 8192)
		multiWriter := io.MultiWriter(os.Stderr, bufferedFile)

		closer := func() error {

			if err := bufferedFile.Flush(); err != nil {
				return fmt.Errorf("failed to flush log file: ", err)
			}
			if err = file.Close(); err != nil {
				return fmt.Errorf("failed to close log file: ", err)
			}
			return nil
		}
		return log.New(multiWriter, "", log.LstdFlags), closer, nil
	}
	close := func() error {
		return nil
	}
	return log.New(os.Stderr, "", log.LstdFlags), close, nil
}
