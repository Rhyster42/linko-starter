package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"

	pkgerr "github.com/pkg/errors"
)

type closeFunc func() error

type stackTracer interface {
	error
	StackTrace() pkgerr.StackTrace
}

func initializeLogger(logFile string) (*slog.Logger, closeFunc, error) {

	debugHandler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level:       slog.LevelDebug,
		ReplaceAttr: replaceAttr,
	})

	if logFile != "" {
		file, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o644)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to open log file: %w", err)
		}
		bufferedFile := bufio.NewWriterSize(file, 8192)
		multiWriter := io.MultiWriter(os.Stderr, bufferedFile)

		closer := func() error {

			if err := bufferedFile.Flush(); err != nil {
				return fmt.Errorf("failed to flush log file: %v", err)
			}
			if err = file.Close(); err != nil {
				return fmt.Errorf("failed to close log file: %v", err)
			}
			return nil
		}
		infoHandler := slog.NewJSONHandler(multiWriter, &slog.HandlerOptions{
			Level:       slog.LevelInfo,
			ReplaceAttr: replaceAttr,
		})

		return slog.New(slog.NewMultiHandler(infoHandler, debugHandler)), closer, nil
	}
	closer := func() error {
		return nil
	}

	infoHandler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level:       slog.LevelInfo,
		ReplaceAttr: replaceAttr,
	})

	return slog.New(slog.NewMultiHandler(infoHandler, debugHandler)), closer, nil
}

func replaceAttr(groups []string, a slog.Attr) slog.Attr {
	if a.Key == "error" {
		err, ok := a.Value.Any().(error)
		if !ok {
			return a
		}

		if stackErr, ok := errors.AsType[stackTracer](err); ok {
			return slog.GroupAttrs("error", slog.Attr{
				Key:   "message",
				Value: slog.StringValue(stackErr.Error()),
			}, slog.Attr{
				Key:   "stack_trace",
				Value: slog.StringValue(fmt.Sprintf("%+v", stackErr.StackTrace())),
			})
		}
		return slog.String("error", fmt.Sprintf("%+v", err))
	}
	return a
}
