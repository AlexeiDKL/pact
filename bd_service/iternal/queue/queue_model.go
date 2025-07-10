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
	mu           sync.Mutex
	Validation   []ValidationItem
	Download     []DownloadItem
	DownloadCh   chan DownloadItem
	ValidationCh chan ValidationItem
}
