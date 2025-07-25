package basedate

import (
	"fmt"
	"strings"
)

func (d *Database) GetVersion() (string, error) {
	// Здесь должна быть логика получения версии базы данных
	// Например, запрос к базе данных или чтение из файла
	// Возвращаем заглушку для примера
	return "1.0.0", nil
}

func (d *Database) GetLatestVersionsByLanguages(lang string) (int, error) {
	query := `SELECT MAX(v.version) AS max_version
			FROM version v
			JOIN language l ON v.language_id = l.id
			WHERE l.short_name = $1;
			`
	rows, err := d.DB.Query(query, lang)
	if err != nil {
		return -1, fmt.Errorf("ошибка запроса: %w", err)
	}
	defer rows.Close()
	var maxVersion int64
	if rows.Next() {
		if err := rows.Scan(&maxVersion); err != nil {
			return -1, fmt.Errorf("ошибка сканирования: %w", err)
		}
	}
	if maxVersion == -1 {
		return -1, fmt.Errorf("нет версий для языка: %s", lang)
	}
	return int(maxVersion), nil
}

func (d *Database) GetFyleTypeByName(typeName string) (int, error) {
	typeName = strings.TrimSpace(typeName)
	typeName = strings.ToLower(typeName)

	query := `
		SELECT id,
       name,
       description
		FROM public.file_type
		WHERE name = $1;
	`
	rows, err := d.DB.Query(query, typeName)
	if err != nil {
		return -1, fmt.Errorf("ошибка запроса: %w", err)
	}
	defer rows.Close()

	var langType int
	if rows.Next() {
		if err := rows.Scan(&langType); err != nil {
			return -1, fmt.Errorf("ошибка сканирования: %w", err)
		}
	}
	if langType == 0 {
		return 0, fmt.Errorf("такого типа файлов не существует: %d", langType)
	}
	if err != nil {
		return 0, err
	}

	return langType, nil
}

func (d *Database) GetLatestVersionsByLanguagesID(languageIDs []int) ([]Version, error) {
	if len(languageIDs) == 0 {
		return nil, nil
	}

	// Формируем плейсхолдеры $1, $2, ..., $n
	placeholders := make([]string, len(languageIDs))
	args := make([]any, len(languageIDs))
	for i, id := range languageIDs {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	query := fmt.Sprintf(`
        SELECT DISTINCT ON (language_id) id, version, pact_id, contents_id, full_text_id, language_id, created_at, updated_at
        FROM version
        WHERE language_id IN (%s)
        ORDER BY language_id, version DESC
    `, strings.Join(placeholders, ", "))

	rows, err := d.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var versions []Version
	for rows.Next() {
		var v Version
		err := rows.Scan(
			&v.Id,
			&v.Version,
			&v.PactId,
			&v.ContentsId,
			&v.FullTextId,
			&v.LanguageId,
			&v.CreatedAt,
			&v.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		versions = append(versions, v)
	}

	return versions, nil
}
