package basedate

import (
	"time"
)

// todo добавить типы файлов
var FileTypeOther = 5
var FileTypeAttachment = 4
var FileTypePact = 1
var FileTypeContents = 2
var FileTypeFullText = 3

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
	Id           int    `json:"id"`
	Checksum     string `json:"checksum"`
	Name         string `json:"name"`
	FileTypeId   int    `json:"file_type_id"`
	LanguageId   int    `json:"language_id"`
	Topic        string `json:"topic"`
	FilePath     string `json:"file_path"`
	DownloadDate string `json:"download_date"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
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
