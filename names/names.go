package names

import (
	"bufio"
	"os"
)

func GetNames(filePathNames string) ([]string, error) {
	names := make([]string, 0)

	file, err := os.Open(filePathNames)
	if err != nil {
		return names, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		name := scanner.Text()
		names = append(names, name)
	}
	return names, err
}
