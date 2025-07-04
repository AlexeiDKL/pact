package queue

func (qm *QueueManager) AddValidation(item ValidationItem) {
	qm.mu.Lock()
	defer qm.mu.Unlock()
	qm.Validation = append(qm.Validation, item)
}

func (qm *QueueManager) AddDownload(item DownloadItem) {
	qm.mu.Lock()
	defer qm.mu.Unlock()
	qm.Download = append(qm.Download, item)
}

func (qm *QueueManager) GetValidationQueue() []ValidationItem {
	qm.mu.Lock()
	defer qm.mu.Unlock()
	return append([]ValidationItem(nil), qm.Validation...) // защитная копия
}

func (qm *QueueManager) GetDownloadQueue() []DownloadItem {
	qm.mu.Lock()
	defer qm.mu.Unlock()
	return append([]DownloadItem(nil), qm.Download...)
}

func (qm *QueueManager) ClearQueues() {
	qm.mu.Lock()
	defer qm.mu.Unlock()
	qm.Validation = nil
	qm.Download = nil
}

func NewQueueManager() *QueueManager {
	return &QueueManager{
		Validation: make([]ValidationItem, 0),
		Download:   make([]DownloadItem, 0),
	}
}
