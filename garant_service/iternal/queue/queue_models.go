package queue

import "sync"

type ValidationItem struct {
	Topic      string
	LanguageID string
	VersionID  string
	FileType   string
}

type DownloadItem struct {
	Topic      string
	LanguageID int
	VersionID  int
	FileType   int
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
	ID           int
	Checksum     string
	Name         string
	FilePath     string
	Topic        string
	LanguageID   int
	VersionID    int
	FileTypeID   int
	DownloadTime string
	CreatedAt    string
	UpdateAt     string
}

type DocumentServiceItem struct {
	Body BDFile
}
