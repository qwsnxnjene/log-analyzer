package analyzer

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func CountByLevel(filename, level string) (int, error) {
	if level != "Error" && level != "Info" && level != "Debug" {
		return 0, fmt.Errorf("[CountByLevel]: unknown log level: %s", level)
	}

	file, err := os.Open(filename)
	if err != nil {
		return 0, fmt.Errorf("[CountByLevel]: error opening file: %v", err)
	}
	defer file.Close()

	buff := bufio.NewScanner(file)
	counter := 0
	for buff.Scan() {
		line := buff.Text()
		if strings.Contains(line, fmt.Sprintf("[%s]", strings.ToUpper(level))) {
			counter++
		}
	}
	if buff.Err() != nil {
		return 0, fmt.Errorf("[CountByLevel]: error reading file: %v", buff.Err())
	}

	return counter, nil
}

func FilterLogs(filename, level, keyword string) ([]string, error) {
	if level != "Error" && level != "Info" && level != "Debug" {
		return nil, fmt.Errorf("[FilterLogs]: unknown log level: %s", level)
	}

	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("[FilterLogs]: error opening file: %v", err)
	}
	defer file.Close()

	buff := bufio.NewScanner(file)
	res := []string{}
	for buff.Scan() {
		line := buff.Text()
		if strings.Contains(line, fmt.Sprintf("[%s]", strings.ToUpper(level))) {
			if strings.Contains(strings.ToLower(line), strings.ToLower(keyword)) {
				res = append(res, line)
			}
		}
	}
	if buff.Err() != nil {
		return nil, fmt.Errorf("[FilterLogs]: error reading file: %v", buff.Err())
	}

	return res, nil
}
