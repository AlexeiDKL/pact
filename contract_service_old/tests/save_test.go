package tests

import (
	"bytes"
	"io"
	"os"
	"testing"

	"dkl.ru/pact/contract_service_old/iternal/files"
)

// TestSaveString проверяет сохранение строки в файл
func TestSaveString(t *testing.T) {
	path := "test_string.txt"
	data := "Hello, Go!"

	err := files.Save(path, data)
	if err != nil {
		t.Fatalf("Ошибка сохранения строки: %s", err)
	}

	// Проверяем, существует ли файл
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatalf("Файл не был создан")
	}

	// Удаляем тестовый файл
	os.Remove(path)
}

// TestSaveByteSlice проверяет сохранение []byte в файл
func TestSaveByteSlice(t *testing.T) {
	path := "test_bytes.txt"
	data := []byte("Byte data")

	err := files.Save(path, data)
	if err != nil {
		t.Fatalf("Ошибка сохранения []byte: %s", err)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatalf("Файл не был создан")
	}

	os.Remove(path)
}

// TestSaveReadCloser проверяет сохранение io.ReadCloser в файл
func TestSaveReadCloser(t *testing.T) {
	path := "test_readcloser.txt"
	data := io.NopCloser(bytes.NewReader([]byte("ReadCloser Data")))

	err := files.Save(path, data)
	if err != nil {
		t.Fatalf("Ошибка сохранения io.ReadCloser: %s", err)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatalf("Файл не был создан")
	}

	os.Remove(path)
}

// TestSaveUnsupportedType проверяет обработку неподдерживаемого типа данных
func TestSaveUnsupportedType(t *testing.T) {
	path := "test_invalid.txt"
	data := 42 // Неподдерживаемый тип

	err := files.Save(path, data)
	if err == nil {
		t.Fatalf("Ожидалась ошибка, но её нет")
	}
}
