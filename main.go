package main

import (
	"flag"
	"fmt"
	"github.com/qwsnxnjene/log-analyzer/analyzer"
	"github.com/qwsnxnjene/log-analyzer/logger"
	"log"
	"os"
	"strings"
)

func main() {
	level := flag.String("level", "", "log level")
	keyword := flag.String("keyword", "", "keyword to filter logs")

	flag.Parse()

	file, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("error creating file: %v", err)
		return
	}

	consoleLogger := logger.NewConsoleLogger()
	fileLogger, err := logger.NewFileLogger(file.Name())
	if err != nil {
		log.Fatalf("error creating FileLogger: %v", err)
	}

	multiLogger := logger.NewMultiLogger([]logger.Logger{consoleLogger, fileLogger}...)
	defer multiLogger.Close()

	command := flag.Args()[0]
	switch command {
	case "log":
		if flag.NArg() < 3 {
			fmt.Println("Usage: log-analyzer log <level> <message>")
			fmt.Println("Example: log-analyzer log error \"Server failed\"")
			os.Exit(1)
		}
		levelLog := flag.Args()[1]
		message := strings.Join(flag.Args()[2:], " ")

		if err := multiLogger.Log(levelLog, message); err != nil {
			log.Fatalf("error logging: %v", err)
		}
	case "analyze":
		if flag.NArg() < 2 {
			fmt.Println("Usage: log-analyzer analyze <filename> --level=<level> --message=<message>")
			fmt.Println("Example: log-analyzer analyze log.txt --level=error --message=\"Server failed\"")
			os.Exit(1)
		}
		filename := flag.Args()[1]
		count, err := analyzer.CountByLevel(filename, *level)
		if err != nil {
			log.Fatalf("error counting: %v", err)
		}
		err = multiLogger.Log("Info", fmt.Sprintf("counted level %s: %d", *level, count))
		if err != nil {
			log.Fatalf("error logging: %v", err)
		}

		lines, err := analyzer.FilterLogs(filename, *level, *keyword)
		if err != nil {
			log.Fatalf("error filtering: %v", err)
		}
		err = multiLogger.Log("Info", fmt.Sprintf("filtered lines %s: %s", *level, lines))
		if err != nil {
			log.Fatalf("error logging: %v", err)
		}
	}
}
