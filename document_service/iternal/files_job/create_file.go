package filesjob

import (
	"os"
)

func CreateFile(path, content string) error {
	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		return err
	}
	return nil
}
