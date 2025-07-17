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
	LanguageID string
	VersionID  string
	FileType   string
}

type QueueManager struct {
	MU                sync.Mutex
	Validation        []ValidationItem
	Download          []DownloadItem
	DocumentService   []DocumentServiceItem
	DocumentServiceCh chan ValidationItem
	DownloadCh        chan DownloadItem
	ValidationCh      chan ValidationItem
}

type DocumentServiceItem struct {
	Topic       string
	LanguageID  string
	FileType    string
	FileVersion string
	FileName    string
}
