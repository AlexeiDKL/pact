package basedate

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestGetLastVersionByLangId(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("ошибка создания mock db: %v", err)
	}
	defer db.Close()

	// Подготавливаем ожидаемый результат
	rows := sqlmock.NewRows([]string{
		"id", "version", "pact_id", "contents_id", "full_text_id", "language_id", "created_at", "updated_at",
	}).AddRow(
		1, 1755613446, 12, sql.NullInt64{}, sql.NullInt64{}, 5,
		time.Date(2025, 8, 19, 17, 24, 6, 0, time.FixedZone("MSK", 3*60*60)),
		time.Time{},
	)

	mock.ExpectQuery("SELECT id,").
		WithArgs(5).
		WillReturnRows(rows)

	// Предполагаемая структура Database с методом BDQueryRow
	d := &Database{DB: db}

	version, err := d.GetLastVersionByLangId(5)
	if err != nil {
		t.Fatalf("ожидали nil ошибку, получили: %v", err)
	}

	if version.Id != 1 ||
		version.Version != 1755613446 ||
		version.PactId != 12 ||
		version.LanguageId != 5 {
		t.Errorf("неожиданные значения version: %+v", version)
	}
}
