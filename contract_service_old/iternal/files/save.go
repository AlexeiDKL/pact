package files

import (
	"fmt"
	"io"
	"os"
)

func Save[T any](path string, data T) error {
	switch v := any(data).(type) {

	case string:
		return newFileString(path, v)
	case io.ReadCloser:
		return newFileFromPost(path, v)
	case []byte:
		return os.WriteFile(path, []byte(v), 0644)
	default:
		return fmt.Errorf("неподдерживаемый тип данных")
	}
}

func newFileFromPost(path string, content io.ReadCloser) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close() // Копирование содержимого в файл
	_, err = io.Copy(file, content)
	return err
}

func newFileString(path string, content string) error { // Создаем файл
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("ошибка при создании файла: %w", err)
	}
	defer file.Close() // Записываем содержимое в файл
	_, err = file.WriteString(content)
	if err != nil {
		return fmt.Errorf("ошибка при записи в файл: %w", err)
	}
	return nil
}
