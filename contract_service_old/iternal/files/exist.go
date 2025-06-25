package files

import "os"

func Exists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err) // Возвращает true, если файл существует
}
