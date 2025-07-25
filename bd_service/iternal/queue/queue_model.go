package queue

import (
	"sync"
)

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

type ValidationItem struct { //todo передаём инфу ту же что и пишем в бд files
	Body BDFile
}

type DownloadItem struct { //todo передаём инфу ту же что и пишем в бд files
	Body BDFile
}

type QueueManager struct {
	MU           sync.Mutex
	Validation   []ValidationItem
	Download     []DownloadItem
	DownloadCh   chan DownloadItem
	ValidationCh chan ValidationItem
}
