package tests

import (
	"os"
	"testing"

	"dkl.ru/pact/contract_service_old/iternal/files"
)

// TestDeleteSuccess проверяет успешное удаление файла
func TestDeleteSuccess(t *testing.T) {
	path := "test_delete.txt"

	// Создаем временный файл
	file, err := os.Create(path)
	if err != nil {
		t.Fatalf("Ошибка создания файла: %s", err)
	}
	file.Close()

	// Удаляем файл
	err = files.Delete(path)
	if err != nil {
		t.Fatalf("Ошибка удаления файла: %s", err)
	}

	// Проверяем, что файл действительно удален
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("Файл не был удалён")
	}
}

// TestDeleteNonExistentFile проверяет удаление несуществующего файла
func TestDeleteNonExistentFile(t *testing.T) {
	path := "nonexistent_file.txt"

	err := files.Delete(path)
	if err == nil {
		t.Fatalf("Ожидалась ошибка, но её нет")
	}
}
