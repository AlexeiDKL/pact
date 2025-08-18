package basedate

import (
	"fmt"
	"strconv"
	"time"

	"github.com/lib/pq"
)

func (d *Database) SaveFile(file File) (int, error) {
	const query = `
        INSERT INTO file (
            checksum, name, file_type_id, language_id, topic, file_path, download_date, created_at, updated_at
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8, $9
        )
    	RETURNING id
    `

	fmt.Println("Проверяем файл:", file.Name, "с контрольной суммой:", file.Checksum)
	fmt.Printf("структура файла: %+v\n", file)

	downloadDate, err := stringToTime(file.DownloadDate)
	if err != nil {
		return -1, fmt.Errorf("не удалось преобразовать download_date: %w", err)
	}
	createdAt, err := stringToTime(file.CreatedAt)
	if err != nil {
		return -1, fmt.Errorf("не удалось преобразовать updated_at: %w", err)
	}
	updatedAt, err := stringToTime(file.UpdatedAt)
	if err != nil {
		return -1, fmt.Errorf("не удалось преобразовать updated_at: %w", err)
	}

	var id int64
	err = d.DB.QueryRow(
		query,
		file.Checksum,
		file.Name,
		file.FileTypeId,
		file.LanguageId,
		file.Topic,
		file.FilePath,
		downloadDate, // должно быть в формате time.Time
		createdAt,    // должно быть в формате time.Time
		updatedAt,    // должно быть в формате time.Time
	).Scan(&id)
	if err != nil {
		// Проверяем, что это ошибка дубликата
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			// Дубликат, ищем id по checksum
			err2 := d.DB.QueryRow("SELECT id FROM file WHERE checksum = $1", file.Checksum).Scan(&id)
			if err2 != nil {
				return -1, fmt.Errorf("дубликат, но не удалось найти id: %w", err2)
			}
			return int(id), nil
		}
		return -1, fmt.Errorf("не удалось сохранить файл: %w", err)
	}
	return int(id), nil
}

func stringToTime(s string) (time.Time, error) {
	if s == "" {
		return time.Time{}, nil
	}
	timestamp, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return time.Time{}, fmt.Errorf("не удалось преобразовать строку в время: %w", err)
	}
	t := time.Unix(timestamp, 0).UTC()
	return t, nil
}
