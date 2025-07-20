package basedate

import "time"

var FileTypeOther = -1
var FileTypeAttachment = -1
var FileTypeContract = -1
var FileTypeContents = -1
var FileTypeFullText = -1

type FileType struct {
	Id          int
	Name        string
	Description *string
}

type Language struct {
	ID          int     `json:"id"`
	FullName    string  `json:"full_name"`
	ShortName   *string `json:"short_name"`
	Description *string `json:"description"`
}

type File struct {
	Id           int       `json:"id"`
	Checksum     string    `json:"checksum"`
	Name         string    `json:"name"`
	FileTypeId   int       `json:"file_type_id"`
	LanguageId   int       `json:"language_id"`
	Topic        *int64    `json:"topic"`
	FilePath     string    `json:"file_path"`
	DownloadDate time.Time `json:"download_date"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Version struct {
	Id         int       `json:"id"`
	Version    int64     `json:"version"`
	PactId     int       `json:"pact_id"`
	ContentsId int       `json:"contents_id"`
	FullTextId int       `json:"full_text_id"`
	LanguageId int       `json:"language_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type VersionAttachment struct {
	VersionId int `json:"version_id"`
	FileId    int `json:"file_id"`
}

type Log struct {
	Id        int       `json:"id"`
	Service   string    `json:"service"`
	ErrorCode int64     `json:"error_code"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}
