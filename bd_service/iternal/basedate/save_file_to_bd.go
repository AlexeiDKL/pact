package basedate

import "fmt"

func (d *Database) SaveFile(versionID int, fileType, filePath, checksum string) error {
	const query = `
        INSERT INTO files (version_id, file_type, file_path, checksum)
        VALUES ($1, $2, $3, $4)
    `
	_, err := d.DB.Exec(query, versionID, fileType, filePath, checksum)
	if err != nil {
		return fmt.Errorf("не удалось сохранить файл: %w", err)
	}
	return nil
}
