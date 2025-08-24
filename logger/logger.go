package logger

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

type Logger interface {
	Log(level, message string) error
	Close() error
}

type FileLogger struct {
	file   *os.File
	writer *bufio.Writer
	ch     chan string
	wg     sync.WaitGroup
}

const DefaultBufferSize = 4096

func NewFileLogger(filename string) (*FileLogger, error) {
	return NewFileLoggerWithBuffer(filename, DefaultBufferSize)
}

func NewFileLoggerWithBuffer(filename string, bufferSize int) (*FileLogger, error) {
	if bufferSize <= 0 {
		bufferSize = DefaultBufferSize
	}

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, fmt.Errorf("[NewFileLogger]: error opening file: %v", err)
	}

	writer := bufio.NewWriterSize(file, bufferSize)
	ch := make(chan string, 100)
	logger := &FileLogger{
		file:   file,
		writer: writer,
		ch:     ch,
	}

	go logger.run()
	logger.wg.Add(1)

	return logger, nil
}

func (f *FileLogger) run() {
	defer f.wg.Done()

	for msg := range f.ch {
		if _, err := f.writer.WriteString(msg + "\n"); err != nil {
			fmt.Printf("[FileLogger.run]: error writing to file: %v\n", err)
			continue
		}
		if err := f.writer.Flush(); err != nil {
			fmt.Printf("[FileLogger.run]: error writing to file: %v\n", err)
		}
	}
}

func (f *FileLogger) Log(level, message string) error {
	if level != "Error" && level != "Info" && level != "Debug" {
		return fmt.Errorf("[f.Log]: unknown log level: %s", level)
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logToWrite := fmt.Sprintf("[%s] [%s] %s", timestamp, strings.ToUpper(level), message)

	select {
	case f.ch <- logToWrite:
		return nil
	default:
		return fmt.Errorf("[f.Log]: channel buffer is full")
	}
}

func (f *FileLogger) Close() error {
	close(f.ch)
	f.wg.Wait()
	return f.file.Close()
}

type ConsoleLogger struct{}

func NewConsoleLogger() *ConsoleLogger {
	return &ConsoleLogger{}
}

func (c *ConsoleLogger) Log(level, message string) error {
	if level != "Error" && level != "Info" && level != "Debug" {
		return fmt.Errorf("[c.Log]: unknown log level: %s", level)
	}
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logToWrite := fmt.Sprintf("[%s] [%s] %s", timestamp, strings.ToUpper(level), message)

	fmt.Println(logToWrite)

	return nil
}

func (c *ConsoleLogger) Close() error {
	return nil
}

type MultiLogger struct {
	loggers []Logger
}

func NewMultiLogger(loggers ...Logger) *MultiLogger {
	return &MultiLogger{loggers: loggers}
}

func (m *MultiLogger) Log(level, message string) error {
	if level != "Error" && level != "Info" && level != "Debug" {
		return fmt.Errorf("[m.Log]: unknown log level: %s", level)
	}

	for _, logger := range m.loggers {
		err := logger.Log(level, message)
		if err != nil {
			return fmt.Errorf("[m.Log]: error logging: %v", err)
		}
	}

	return nil
}

func (m *MultiLogger) Close() error {
	for _, logger := range m.loggers {
		err := logger.Close()
		if err != nil {
			return fmt.Errorf("[m.Log]: error closing logger: %v", err)
		}
	}

	return nil
}
