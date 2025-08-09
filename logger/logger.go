package logger

import (
	"bufio"
	"fmt"
	"os"
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

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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
	logToWrite := fmt.Sprintf("[%s] [%s] %s", timestamp, level, message)

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
