package queue

import (
	"sync"
)

type BDFile struct {
	ID           int    `json:"id"`
	Checksum     string `json:"checksum"`
	Name         string `json:"name"`
	FilePath     string `json:"file_path"`
	Topic        string `json:"topic"`
	LanguageID   int    `json:"language_id"`
	VersionID    int
	FileTypeID   int    `json:"file_type_id"`
	DownloadTime string `json:"download_date"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

type ValidationItem struct { //todo передаём инфу ту же что и пишем в бд files
	Body BDFile
}

type DownloadItem struct { //todo передаём инфу ту же что и пишем в бд files
	Body BDFile
}

type QueueManager struct {
	MU         sync.Mutex
	Validation []ValidationItem
	Download   []DownloadItem
	Version    []BDFile // Список версий, которые нужно сохранить

	VersionCh    chan BDFile
	DownloadCh   chan DownloadItem
	ValidationCh chan ValidationItem
}
