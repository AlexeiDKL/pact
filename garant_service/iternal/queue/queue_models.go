package queue

import (
	"sync"
)

type ValidationRequest struct {
	Body FileItem `json:"Body"`
}

type DownloadItem struct {
	Body FileItem `json:"Body"`
}

type ValidationItem struct {
	Body FileItem `json:"Body"`
}

type FileItem struct {
	Topic      string `json:"topic"`
	LanguageID int    `json:"language_id"`
	VersionID  int    `json:"version_id"`
	FileType   int    `json:"file_type"`
}

type QueueManager struct {
	MU              sync.Mutex
	Validation      []ValidationItem
	Download        []DownloadItem
	DocumentService []DocumentServiceItem
	SaveBdFile      []BDFile

	DownloadCh        chan DownloadItem
	ValidationCh      chan ValidationItem
	DocumentServiceCh chan DocumentServiceItem
	SaveBDFileCH      chan BDFile
}

type BDFile struct {
	ID           int    `json:"id"`
	Checksum     string `json:"checksum"`
	Name         string `json:"name"`
	FilePath     string `json:"file_path"`
	Topic        string `json:"topic"`
	LanguageID   int    `json:"language_id"`
	VersionID    int    `json:"version_id"`
	FileTypeID   int    `json:"file_type_id"`
	DownloadTime string `json:"download_date"`
	CreatedAt    string `json:"created_at"`
	UpdateAt     string `json:"updated_at"`
}

type DocumentServiceItem struct {
	Body BDFile
}
