package basedate_test

import (
	"testing"

	"dkl.ru/pact/bd_service/iternal/basedate"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestSaveFile(t *testing.T) {
	// 1. Создаём мок соединения с БД
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// 2. Инициализируем Database
	database := &basedate.Database{DB: db}

	// 3. Ожидаем SQL и параметры
	mock.ExpectExec(`INSERT INTO files`).
		WithArgs(1, "contract", "/path/to/file.odt", "abc123").
		WillReturnResult(sqlmock.NewResult(1, 1))

	// 4. Вызываем целевую функцию
	err = database.SaveFile(1, "contract", "/path/to/file.odt", "abc123")
	assert.NoError(t, err)

	// 5. Проверяем, что все ожидания mock выполнены
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
