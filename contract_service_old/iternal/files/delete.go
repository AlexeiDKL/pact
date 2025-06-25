package files

import (
	"fmt"
	"os"
)

func Delete(filePath string) error {
	err := os.Remove(filePath)
	if err != nil {
		return fmt.Errorf("ошибка удаления файла %s: %w", filePath, err)
	}
	return nil
}
