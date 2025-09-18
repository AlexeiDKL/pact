package basedate

import (
	"database/sql"
	"fmt"
)

func (d *Database) BdExec(query string, args ...interface{}) error {

	// Здесь будет логика выполнения базовых операций с данными
	// Например, инициализация базы данных, выполнение запросов и т.д.
	_, err := d.DB.Exec(query, args...)
	return err
}

func (d *Database) BDQueryRow(query string, args ...interface{}) *sql.Row {
	return d.DB.QueryRow(query, args...)
}

func (d *Database) BDQuery(query string, args ...interface{}) (*sql.Rows, error) {
	// Здесь будет логика выполнения запроса к базе данных
	// Например, выполнение SELECT запроса и возврат множества строк результата
	rows, err := d.DB.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err)
	}
	return rows, nil
}
