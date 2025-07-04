package queue

import "sync"

type ValidationItem struct {
	Topic      string
	LanguageID int
	VersionID  int64
	FileType   string
}

type DownloadItem struct {
	Topic      string
	LanguageID int
	Reason     string
}

type QueueManager struct {
	mu         sync.Mutex
	Validation []ValidationItem
	Download   []DownloadItem
}
