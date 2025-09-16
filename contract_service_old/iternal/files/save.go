package files

import (
	"fmt"
	"io"
	"os"

	"dkl.ru/pact/contract_service_old/iternal/toc"

	"encoding/json"
)

func Save[T any](path string, data T) error {
	switch v := any(data).(type) {
	case string:
		return newFileString(path, v)
	case io.ReadCloser:
		return newFileFromPost(path, v)
	case []byte:
		return os.WriteFile(path, []byte(v), 0644)
	case []toc.TOCItem:
		return newFileToc(path, v)
	default:
		// Попробовать сериализовать через json.Marshal
		file, err := os.Create(path)
		if err != nil {
			return fmt.Errorf("ошибка при создании файла: %w", err)
		}
		defer file.Close()
		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(data); err != nil {
			return fmt.Errorf("ошибка при кодировании JSON: %w", err)
		}
		return nil
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

func newFileToc(path string, content []toc.TOCItem) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("ошибка при создании файла: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(content); err != nil {
		return fmt.Errorf("ошибка при кодировании JSON: %w", err)
	}
	return nil
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

func SaveToJSON(path string, data any) error {
	return Save(path, data)
}
