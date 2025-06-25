package basedate

import (
	"database/sql"
	"fmt"

	"dkl.ru/pact/bd_service/iternal/logger"
)

func valideTable(db *sql.DB) (bool, error) {
	schemas := map[string][]ColumnInfo{
		"language": {
			{"id", "integer", "NO"},
			{"name", "character varying", "NO"},
		},
		"versions": {
			{"id", "integer", "NO"},
			{"version_number", "character varying", "NO"},
			{"language_id", "bigint", "NO"},
			{"created_at", "timestamp with time zone", "NO"},
			{"size", "bigint", "NO"},
			{"is_latest", "boolean", "NO"},
		},
		"files": {
			{"id", "integer", "NO"},
			{"version_id", "bigint", "NO"},
			{"file_type", "character varying", "NO"},
			{"file_path", "character varying", "NO"},
			{"checksum", "character varying", "NO"},
		},
		"fulltext": {
			{"id", "integer", "NO"},
			{"version_id", "bigint", "NO"},
			{"content", "character varying", "NO"},
		},
		"contents": {
			{"id", "integer", "NO"},
			{"version_id", "bigint", "NO"},
			{"content", "character varying", "NO"},
		},
		"logs": {
			{"id", "integer", "NO"},
			{"service", "character varying", "NO"},
			{"error_code", "bigint", "NO"},
			{"message", "character varying", "NO"},
			{"created_at", "timestamp with time zone", "NO"},
		},
		"cache": {
			{"value", "character varying", "NO"},
			{"expires_at", "timestamp with time zone", "NO"},
		},
	}

	for tableName, expectedCols := range schemas {
		err := ValidateTableSchema(db, tableName, expectedCols)
		if err != nil {
			return false, fmt.Errorf("ошибка в таблице %s: %v", tableName, err)
		} else {
			logger.Logger.Debug(fmt.Sprintf("✅ Структура таблицы %s корректна", tableName))
		}
	}
	return true, nil
}

type ColumnInfo struct {
	ColumnName string
	DataType   string
	IsNullable string
}

func ValidateTableSchema(db *sql.DB, tableName string, expectedCols []ColumnInfo) error {
	rows, err := db.Query(`
        SELECT column_name, data_type, is_nullable
        FROM information_schema.columns
        WHERE table_name = $1
        ORDER BY ordinal_position`, tableName)
	if err != nil {
		return fmt.Errorf("ошибка запроса к information_schema: %w", err)
	}
	defer rows.Close()

	var actualCols []ColumnInfo
	for rows.Next() {
		var col ColumnInfo
		if err := rows.Scan(&col.ColumnName, &col.DataType, &col.IsNullable); err != nil {
			return fmt.Errorf("ошибка чтения столбцов: %w", err)
		}
		actualCols = append(actualCols, col)
	}

	if len(actualCols) != len(expectedCols) {
		return fmt.Errorf("❌ %s: ожидалось %d столбцов, найдено %d", tableName, len(expectedCols), len(actualCols))
	}

	for i, expected := range expectedCols {
		actual := actualCols[i]
		if actual != expected {
			return fmt.Errorf("❌ %s: столбец %d отличается:\n  ожидалось: %+v\n  фактически: %+v", tableName, i+1, expected, actual)
		}
	}

	return nil
}
