package basedate

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

func ExecuteDropTablesSQL(db *sql.DB, filepath string) error {
	sqlBytes, err := os.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("ошибка при чтении SQL-файла: %w", err)
	}

	sqlContent := string(sqlBytes)

	// Разделяем SQL-команды по точке с запятой
	commands := strings.Split(sqlContent, ";")

	for _, cmd := range commands {
		cmd = strings.TrimSpace(cmd)
		if cmd == "" {
			continue
		}
		_, err := db.Exec(cmd)
		if err != nil {
			return fmt.Errorf("ошибка при выполнении команды [%s]: %w", cmd, err)
		}
	}

	return nil
}
