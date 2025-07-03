package basedate_test

import (
	"testing"

	"dkl.ru/pact/bd_service/iternal/basedate"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func setupMockDB(t *testing.T) (*basedate.Database, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	return &basedate.Database{DB: db}, mock
}

func TestSaveFile(t *testing.T) {
	// 1. Создаём мок соединения с БД
	database, mock := setupMockDB(t)

	// 3. Ожидаем SQL и параметры
	mock.ExpectExec(`INSERT INTO files`).
		WithArgs(1, "contract", "/path/to/file.odt", "abc123").
		WillReturnResult(sqlmock.NewResult(1, 1))

	// 4. Вызываем целевую функцию
	var q string = "abc123" // Пример контрольной суммы
	err := database.SaveFile(1, "contract", "/path/to/file.odt", q)
	assert.NoError(t, err)

	// 5. Проверяем, что все ожидания mock выполнены
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
