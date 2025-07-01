package tests

import (
	"os"
	"testing"

	"dkl.ru/pact/contract_service_old/iternal/files"
)

func TestExists(t *testing.T) {
	path := "test_exists.txt"

	// Создаем временный файл
	file, err := os.Create(path)
	if err != nil {
		t.Fatalf("Ошибка создания файла: %s", err)
	}
	file.Close()

	// Проверяем, что файл существует
	if !files.Exists(path) {
		t.Fatalf("Файл должен существовать, но функция вернула false")
	}

	// Удаляем файл
	os.Remove(path)

	// Проверяем, что файл больше не существует
	if files.Exists(path) {
		t.Fatalf("Файл был удалён, но функция вернула true")
	}
}
