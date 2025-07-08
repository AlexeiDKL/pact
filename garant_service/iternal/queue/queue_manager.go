package queue

func NewQueueManager() *QueueManager {
	return &QueueManager{
		Validation: make([]ValidationItem, 0),
		Download:   make([]DownloadItem, 0),
	}
}

func (qm *QueueManager) AddValidation(item ValidationItem) {
	qm.MU.Lock()
	defer qm.MU.Unlock()
	qm.Validation = append(qm.Validation, item)
}

func (qm *QueueManager) AddDownload(item DownloadItem) {
	qm.MU.Lock()
	defer qm.MU.Unlock()
	qm.Download = append(qm.Download, item)
}

func (qm *QueueManager) GetValidationQueue() []ValidationItem {
	qm.MU.Lock()
	defer qm.MU.Unlock()
	return append([]ValidationItem(nil), qm.Validation...) // защитная копия
}

func (qm *QueueManager) GetDownloadQueue() []DownloadItem {
	qm.MU.Lock()
	defer qm.MU.Unlock()
	return append([]DownloadItem(nil), qm.Download...)
}
