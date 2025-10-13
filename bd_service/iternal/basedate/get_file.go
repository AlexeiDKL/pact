package basedate

import (
	"context"
	"database/sql"
	"fmt"
)

func (d *Database) GetFileNameByID(ctx context.Context, fileID int) (string, error) {
	var fileName string
	query := `SELECT file_name FROM file WHERE id = $1`
	if err := d.DB.QueryRowContext(ctx, query, fileID).Scan(&fileName); err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("файл с ID %d не найден", fileID)
		}
		return "", fmt.Errorf("ошибка запроса имени файла по ID %d: %w", fileID, err)
	}
	return fileName, nil
}
