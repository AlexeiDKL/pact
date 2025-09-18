package basedate

func (d *Database) GetLastVersionByLangId(langId int) (Version, error) {
	query := `
	SELECT id,
       version,
       pact_id,
       contents_id,
       full_text_id,
       language_id,
       created_at,
       updated_at
	FROM public.version
	WHERE language_id = $1
	ORDER BY version DESC
	LIMIT 1;
	`
	var version Version
	row := d.BDQueryRow(query, langId)

	err := row.Scan(&version.Id,
		&version.Version,
		&version.PactId,
		&version.ContentsId,
		&version.FullTextId,
		&version.LanguageId,
		&version.CreatedAt,
		&version.UpdatedAt,
	)
	return version, err

}
