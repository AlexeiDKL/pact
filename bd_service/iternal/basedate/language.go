package basedate

import "strings"

func (d *Database) GetLanguageWithID(languageID int) (Language, error) {
	const query = `
		SELECT id, full_name, short_name, description
		FROM language
		WHERE id = $1
	`
	var lang Language
	err := d.DB.QueryRow(query, languageID).Scan(&lang.ID, &lang.FullName, &lang.ShortName, &lang.Description)
	if err != nil {
		return Language{}, err
	}
	return lang, nil
}

func (d *Database) GetAllLanguages() ([]Language, error) {
	const query = `
		SELECT id, full_name, short_name, description
		FROM language
	`
	rows, err := d.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var languages []Language
	for rows.Next() {
		var lang Language
		if err := rows.Scan(&lang.ID, &lang.FullName, &lang.ShortName, &lang.Description); err != nil {
			return nil, err
		}
		languages = append(languages, lang)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return languages, nil
}

func (d *Database) SaveLanguage(lang Language) error {
	const query = `
		INSERT INTO language (full_name, short_name, description)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	err := d.DB.QueryRow(query, lang.FullName, lang.ShortName, lang.Description).Scan(&lang.ID)
	if err != nil {
		return err
	}
	return nil
}

func (d *Database) UpdateLanguage(lang Language) error {
	const query = `
		UPDATE language
		SET full_name = $1, short_name = $2, description = $3
		WHERE id = $4
	`
	_, err := d.DB.Exec(query, lang.FullName, lang.ShortName, lang.Description, lang.ID)
	if err != nil {
		return err
	}
	return nil
}

func (d *Database) DeleteLanguage(languageID int) error {
	const query = `
		DELETE FROM language
		WHERE id = $1
	`
	_, err := d.DB.Exec(query, languageID)
	if err != nil {
		return err
	}
	return nil
}

func (d *Database) GetLanguageByShortName(shortName string) (Language, error) {
	shortName = strings.TrimSpace(shortName)
	if shortName == "" {
		return Language{}, nil
	}
	shortName = strings.ToLower(shortName)
	const query = `
		SELECT id, full_name, short_name, description
		FROM language
		WHERE short_name = $1
	`
	var lang Language
	err := d.DB.QueryRow(query, shortName).Scan(&lang.ID, &lang.FullName, &lang.ShortName, &lang.Description)
	if err != nil {
		return Language{}, err
	}
	return lang, nil
}

func (d *Database) GetLanguageByFullName(fullName string) (Language, error) {
	fullName = strings.TrimSpace(fullName)
	if fullName == "" {
		return Language{}, nil
	}
	fullName = strings.ToLower(fullName)
	const query = `
		SELECT id, full_name, short_name, description
		FROM language
		WHERE full_name = $1
	`
	var lang Language
	err := d.DB.QueryRow(query, fullName).Scan(&lang.ID, &lang.FullName, &lang.ShortName, &lang.Description)
	if err != nil {
		return Language{}, err
	}
	return lang, nil
}

func (d *Database) GetLanguageByID(languageID int) (Language, error) {
	const query = `
		SELECT id, full_name, short_name, description
		FROM language
		WHERE id = $1
	`
	var lang Language
	err := d.DB.QueryRow(query, languageID).Scan(&lang.ID, &lang.FullName, &lang.ShortName, &lang.Description)
	if err != nil {
		return Language{}, err
	}
	return lang, nil
}

func (d *Database) GetLanguageByName(name string) (Language, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return Language{}, nil
	}
	name = strings.ToLower(name)

	const query = `
		SELECT id, full_name, short_name, description
		FROM language
		WHERE full_name = $1 OR short_name = $1
	`
	var lang Language
	err := d.DB.QueryRow(query, name).Scan(&lang.ID, &lang.FullName, &lang.ShortName, &lang.Description)
	if err != nil {
		return Language{}, err
	}
	return lang, nil
}
