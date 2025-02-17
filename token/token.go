package token

import (
	"io"
	"os"
	"strings"
)

func GetToken(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	data := make([]byte, 64)

	n, err := file.Read(data)
	if err == io.EOF {
		err = nil
	}
	if err != nil {
		return "", err
	}

	token := string(data[:n])
	token = strings.TrimSpace(token)
	return token, err
}
