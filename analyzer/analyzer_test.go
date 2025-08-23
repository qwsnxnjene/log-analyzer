package analyzer

import (
	logger2 "github.com/qwsnxnjene/log-analyzer/logger"
	"os"
	"strings"
	"testing"
)

func TestCountByLevelBasicCount(t *testing.T) {
	file, err := os.CreateTemp("", "testlog-*.txt")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(file.Name())

	logger, err := logger2.NewFileLogger(file.Name())
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}

	err = logger.Log("Error", "error")
	if err != nil {
		t.Errorf("failed to log: %v", err)
	}
	err = logger.Log("Error", "error")
	if err != nil {
		t.Errorf("failed to log: %v", err)
	}
	err = logger.Log("Info", "info")
	if err != nil {
		t.Errorf("failed to log: %v", err)
	}
	err = logger.Log("Info", "info")
	if err != nil {
		t.Errorf("failed to log: %v", err)
	}
	err = logger.Log("Debug", "debug")
	if err != nil {
		t.Errorf("failed to log: %v", err)
	}
	err = logger.Log("Debug", "debug")
	if err != nil {
		t.Errorf("failed to log: %v", err)
	}

	logger.Close()

	if count, err := CountByLevel(file.Name(), "Error"); count != 2 {
		t.Errorf("wrong Error count, expected 2, got %d", count)
	} else if err != nil {
		t.Errorf("failed counting: %v", err)
	}

	if count, err := CountByLevel(file.Name(), "Info"); count != 2 {
		t.Errorf("wrong Info count, expected 2, got %d", count)
	} else if err != nil {
		t.Errorf("failed counting: %v", err)
	}

	if count, err := CountByLevel(file.Name(), "Debug"); count != 2 {
		t.Errorf("wrong Debug count, expected 2, got %d", count)
	} else if err != nil {
		t.Errorf("failed counting: %v", err)
	}
}

func TestCountByLevelNoCount(t *testing.T) {
	file, err := os.CreateTemp("", "testlog-*.txt")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(file.Name())

	logger, err := logger2.NewFileLogger(file.Name())
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}
	defer logger.Close()

	err = logger.Log("Error", "error")
	if err != nil {
		t.Errorf("failed to log: %v", err)
	}

	if count, err := CountByLevel(file.Name(), "Debug"); count != 0 {
		t.Errorf("wrong Debug count, expected 0, got %d", count)
	} else if err != nil {
		t.Errorf("failed counting: %v", err)
	}
}

func TestCountByLevelWrongLevel(t *testing.T) {
	file, err := os.CreateTemp("", "testlog-*.txt")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(file.Name())

	if _, err := CountByLevel(file.Name(), "Bro"); err == nil {
		t.Fatalf("fail to determine wrong log level")
	}
}

func TestFilterLogsBasic(t *testing.T) {
	file, err := os.CreateTemp("", "testlog-*.txt")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(file.Name())

	logger, err := logger2.NewFileLogger(file.Name())
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}

	err = logger.Log("Error", "bruh")
	if err != nil {
		t.Errorf("failed to log: %v", err)
	}
	err = logger.Log("Error", "no")
	if err != nil {
		t.Errorf("failed to log: %v", err)
	}
	err = logger.Log("Info", "bruh")
	if err != nil {
		t.Errorf("failed to log: %v", err)
	}
	err = logger.Log("Info", "goal")
	if err != nil {
		t.Errorf("failed to log: %v", err)
	}
	err = logger.Log("Debug", "goal")
	if err != nil {
		t.Errorf("failed to log: %v", err)
	}
	err = logger.Log("Debug", "no")
	if err != nil {
		t.Errorf("failed to log: %v", err)
	}

	logger.Close()

	if lines, err := FilterLogs(file.Name(), "Info", "bruh"); err == nil {
		if len(lines) != 1 {
			t.Errorf("wrong lines count, expected 1, got %d", len(lines))
		} else {
			if strings.Split(lines[0], "] ")[2] != "bruh" {
				t.Errorf("wrong line filtered: %s", lines)
			}
		}
	} else {
		t.Errorf("error filtering: %v", err)
	}
}
