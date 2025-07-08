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
	Reason     string
}

type QueueManager struct {
	MU         sync.Mutex
	Validation []ValidationItem
	Download   []DownloadItem
}
