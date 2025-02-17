package notification

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func ExcludeLines(filePath string, needToDelete string) (bool, error) {
	lineFound := false
	file, err := os.Open(filePath)
	if err != nil {
		return false, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		text := scanner.Text()
		if strings.Contains(text, needToDelete) {
			lineFound = true
			continue
		}
		lines = append(lines, text)
	}

	if err := scanner.Err(); err != nil {
		return lineFound, err
	}

	err = ioutil.WriteFile(filePath, []byte(strings.Join(lines, "\n")), 0644)
	return lineFound, err
}

func AppendLine(filePath string, needToAppend string) error {
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}

	defer f.Close()

	_, err = f.WriteString(fmt.Sprint("\n", needToAppend))
	return err
}
