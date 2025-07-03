package basedate

func (d *Database) GetFileListWithVersion(versionID int) ([]File, error) {
	const query = `
		SELECT f.id, f.file_type, f.file_path, f.checksum
		FROM files f
		JOIN versions v ON f.version_id = v.id
		WHERE v.id = $1
	`

	rows, err := d.DB.Query(query, versionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []File
	for rows.Next() {
		var file File
		if err := rows.Scan(&file.ID, &file.FileType, &file.FilePath, &file.Checksum); err != nil {
			return nil, err
		}
		files = append(files, file)
	}

	return files, nil
}
