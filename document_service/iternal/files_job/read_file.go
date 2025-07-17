package filesjob

import (
	"fmt"
	"os"
)

func ReadFileToString(filePath string) (string, error) {
	res, err := ReadFileToBytes(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file %s: %w", filePath, err)
	}
	return string(res), nil
}

func ReadFileToBytes(filePath string) ([]byte, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}
	return content, nil
}
