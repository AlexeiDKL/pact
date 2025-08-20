package basedate

import (
	"fmt"
	"strings"
)

func InitializeFileTypes(d *Database) error {
	// ходим в бд и получаем типы файлов
	// инициализируем глобальные переменные
	// типы: договор, приложение, оглавление, полный текст, другие
	// имя таблицы: file_type
	const query = `
		SELECT id, name
		FROM file_type
	`
	rows, err := d.BDQuery(query)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var fileType FileType
		if err := rows.Scan(&fileType.Id, &fileType.Name); err != nil {
			return err
		}
		switch strings.ToLower(fileType.Name) {
		case "pact":
			FileTypePact = fileType.Id
		case "attachment":
			FileTypeAttachment = fileType.Id
		case "contents":
			FileTypeContents = fileType.Id
		case "full_text":
			FileTypeFullText = fileType.Id
		case "other":
			FileTypeOther = fileType.Id
		default:
			return fmt.Errorf("неизвестный тип файла: %s", fileType.Name)
		}
	}
	if err := rows.Err(); err != nil {
		return err
	}
	return nil

}
