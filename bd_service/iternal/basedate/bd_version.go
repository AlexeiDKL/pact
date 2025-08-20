package basedate

import (
	"fmt"
)

func (db *Database) SaveToVersion(file File) error {

	fmt.Printf("Сохранение файла в версию: %+v\n", file)

	switch file.FileTypeId {
	case FileTypePact:
		return saveToVersion(db, file)
	case FileTypeAttachment, FileTypeOther, FileTypeFullText, FileTypeContents:
		{
			var versionId int
			err := db.DB.QueryRow(
				"SELECT id FROM version WHERE language_id = $1 ORDER BY version DESC LIMIT 1",
				file.LanguageId,
			).Scan(&versionId)
			if err != nil {
				return fmt.Errorf("не удалось получить последнюю версию для языка %d: %w", file.LanguageId, err)
			}
			return updateToVersion(db, file, versionId)
		}
	default:
		return fmt.Errorf("неизвестный тип файла: %d", file.FileTypeId)
	}
}

func saveToVersion(d *Database, file File) error {
	// Логика сохранения файла для типа FileTypeContract
	// Здесь должна быть реализация создания новой версии на основе файла
	// version, pact_id, language_id, created_at, updated_at должны быть заполнены в file
	query := `
		INSERT INTO version (version, pact_id, language_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	vsersion := creatVersion(file.CreatedAt)
	createdAt, err := stringToTime(file.CreatedAt)
	if err != nil {
		return fmt.Errorf("не удалось преобразовать created_at: %w", err)
	}
	updateAt, err := stringToTime(file.UpdatedAt)
	if err != nil {
		return fmt.Errorf("не удалось преобразовать updated_at: %w", err)
	}

	fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
	fmt.Printf("Сохраняем версию: %s, pact_id: %d, language_id: %d, created_at: %s, updated_at: %s\n",
		vsersion, file.Id, file.LanguageId, createdAt, updateAt)
	fmt.Printf("Запрос: %s\n", query)
	fmt.Printf("Параметры: %s, %d, %d, %s, %s\n",
		vsersion, file.Id, file.LanguageId, createdAt, updateAt)
	fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")

	errBd := d.DB.QueryRow(
		query,
		vsersion,
		file.Id,         // int или *int
		file.LanguageId, // int
		createdAt,       // time.Time
		updateAt,        // time.Time
	)
	if errBd != nil {
		return fmt.Errorf("не удалось добавить версию: %v", err)
	}
	return nil
}

func creatVersion(createdAt string) string {
	return createdAt
}

func updateToVersion(d *Database, file File, versionId int) error {
	// Логика обновления файла для типов FileTypeAttachment, FileTypeOther, FileTypeFullText, FileTypeContents
	// Здесь должна быть реализация обновления версии, по файлу
	var query string
	switch file.FileTypeId {
	case FileTypeAttachment:
		// если это тип Attachment, то сохраняем их в отдельную таблицу, и добавляем в version айдишки

		query = "INSERT INTO attachment (file_id, version_id) VALUES ($1, $2)"
		_, err := d.DB.Exec(query, file.Id, versionId)
		return err
	case FileTypeFullText:
		// если это тип FullText, то добавляем его id в версию
		query = "UPDATE version SET full_text_id = $1 WHERE id = $2"
		_, err := d.DB.Exec(query, file.Id, versionId)
		return err
	case FileTypeContents:
		query = "UPDATE version SET contents_id = $1 WHERE id = $2"
		_, err := d.DB.Exec(query, file.Id, versionId)
		// если это тип Contents, то добавляем его id в версию
		return err
	default:
		return fmt.Errorf("неизвестный тип файла: %d", file.FileTypeId)
	}
}
