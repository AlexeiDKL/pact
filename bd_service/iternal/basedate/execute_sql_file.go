package basedate

import (
	"database/sql"
	"fmt"
	"io/ioutil"

	_ "github.com/lib/pq"
)

func ExecuteSQLFile(db *sql.DB, filePath string) error {
	sqlBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("ошибка чтения SQL-файла: %v", err)
	}

	sqlScript := string(sqlBytes)
	_, err = db.Exec(sqlScript)
	if err != nil {
		return fmt.Errorf("ошибка выполнения SQL: %v", err)
	}

	return nil
}
