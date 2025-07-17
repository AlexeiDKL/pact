package basedate

import "fmt"

func (d *Database) SaveFile(file File) error {
	const query = `
        INSERT INTO file (
            checksum, name, file_type_id, language_id, topic, file_path, download_date, created_at, updated_at
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8, $9
        )
    `
	_, err := d.DB.Exec(
		query,
		file.Checksum,
		file.Name,
		file.FileTypeId,
		file.LanguageId,
		file.Topic,
		file.FilePath,
		file.DownloadDate,
		file.CreatedAt,
		file.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("не удалось сохранить файл: %w", err)
	}
	return nil
}
