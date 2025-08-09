package logger

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestNewFileLogger(t *testing.T) {
	file, err := os.CreateTemp("", "testlog-*.txt")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(file.Name())

	logger, err := NewFileLogger(file.Name())
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}
	defer logger.Close()

	newFilename := "newtestlog.txt"
	loggerNew, err := NewFileLogger(newFilename)
	if err != nil {
		t.Fatalf("failed to create logger %v: %v", newFilename, err)
	}
	defer loggerNew.Close()
	defer os.Remove(newFilename)
}

func TestNewFileLoggerWithBuffer(t *testing.T) {
	file, err := os.CreateTemp("", "testlog-*.txt")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(file.Name())

	logger, err := NewFileLoggerWithBuffer(file.Name(), 1024)
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}
	defer logger.Close()
	if logger.writer.Size() != 1024 {
		t.Errorf("expected buffer size: %d, got: %d", 1024, logger.writer.Size())
	}

	file, err = os.CreateTemp("", "testlog1-*.txt")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(file.Name())

	loggerZero, err := NewFileLoggerWithBuffer(file.Name(), -10)
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}
	defer loggerZero.Close()
	if loggerZero.writer.Size() != DefaultBufferSize {
		t.Errorf("expected buffer size: %d, got: %d", DefaultBufferSize, logger.writer.Size())
	}
}

func TestFileLogger_Log(t *testing.T) {
	file, err := os.CreateTemp("", "testlog-*.txt")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(file.Name())

	logger, err := NewFileLogger(file.Name())
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}
	defer logger.Close()

	err = logger.Log("Error", "error")
	if err != nil {
		t.Errorf("failed to log: %v", err)
	}

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		text := scanner.Text()
		expected := "[Error] error"
		if !strings.Contains(text, expected) {
			t.Errorf("wrong log format, expected %s - got %s", expected, text)
		}
	}
	if scanner.Err() != nil {
		t.Fatalf("failed to scan file: %v", scanner.Err())
	}

	err = logger.Log("Beta", "b")
	if err != nil && err.Error() != fmt.Errorf("[f.Log]: unknown log level: Beta").Error() {
		t.Errorf("expected error with unknown log level, got %v", err)
	}
}

func TestFileLogger_Close(t *testing.T) {
	file, err := os.CreateTemp("", "testlog-*.txt")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(file.Name())

	logger, err := NewFileLogger(file.Name())
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}

	err = logger.Close()
	if err != nil {
		t.Errorf("failed to close logger: %v", err)
	}
}
