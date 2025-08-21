package queue

func NewQueueManager() *QueueManager {
	return &QueueManager{
		Validation:      make([]ValidationItem, 0),
		Download:        make([]DownloadItem, 0),
		DocumentService: make([]DocumentServiceItem, 0),
		SaveBdFile:      make([]BDFile, 0),

		DownloadCh:        make(chan DownloadItem, 100),
		ValidationCh:      make(chan ValidationItem, 100),
		DocumentServiceCh: make(chan DocumentServiceItem, 100),
		SaveBDFileCH:      make(chan BDFile, 100),
	}
}

func (qm *QueueManager) AddSaveBdFile(item BDFile) {
	qm.MU.Lock()
	defer qm.MU.Unlock()
	qm.SaveBdFile = append(qm.SaveBdFile, item)
	qm.SaveBDFileCH <- item
}

func (qm *QueueManager) RemoveSaveBdFile(target BDFile) {
	qm.MU.Lock()
	defer qm.MU.Unlock()

	newQueue := make([]BDFile, 0, len(qm.Download))
	for _, item := range qm.SaveBdFile {
		if item.FilePath == target.FilePath {
			continue // Пропускаем совпадающий
		}
		newQueue = append(newQueue, item)
	}
	qm.SaveBdFile = newQueue
}

func (qm *QueueManager) AddDocumentService(item DocumentServiceItem) {
	qm.MU.Lock()
	defer qm.MU.Unlock()
	qm.DocumentService = append(qm.DocumentService, item)
	qm.DocumentServiceCh <- item
}

func (qm *QueueManager) AddValidation(item ValidationItem) {
	qm.MU.Lock()
	defer qm.MU.Unlock()
	qm.Validation = append(qm.Validation, item)
	qm.ValidationCh <- item
}

func (qm *QueueManager) AddDownload(item DownloadItem) {
	qm.MU.Lock()
	defer qm.MU.Unlock()
	qm.Download = append(qm.Download, item)
	qm.DownloadCh <- item

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

func (qm *QueueManager) RemoveDocumentServiceItem(target DocumentServiceItem) {
	qm.MU.Lock()
	defer qm.MU.Unlock()

	newQueue := make([]DocumentServiceItem, 0, len(qm.Download))
	for _, item := range qm.DocumentService {
		if item.Body.FilePath == target.Body.FilePath {
			continue // Пропускаем совпадающий
		}
		newQueue = append(newQueue, item)
	}
	qm.DocumentService = newQueue
}

func (qm *QueueManager) RemoveDownloadItem(target DownloadItem) {
	qm.MU.Lock()
	defer qm.MU.Unlock()

	newQueue := make([]DownloadItem, 0, len(qm.Download))
	for _, item := range qm.Download {
		if item.Body.Topic == target.Body.Topic && item.Body.LanguageID == target.Body.LanguageID && item.Body.FileType == target.Body.FileType {
			continue // Пропускаем совпадающий
		}
		newQueue = append(newQueue, item)
	}
	qm.Download = newQueue
}

func (qm *QueueManager) RemoveValidationItem(target ValidationItem) {
	qm.MU.Lock()
	defer qm.MU.Unlock()

	newQueue := make([]ValidationItem, 0, len(qm.Validation))
	for _, item := range qm.Validation {
		if item.Body.Topic == target.Body.Topic && item.Body.LanguageID == target.Body.LanguageID && item.Body.FileType == target.Body.FileType {
			continue
		}
		newQueue = append(newQueue, item)
	}
	qm.Validation = newQueue
}

func (qm *QueueManager) SendFileToDocumentService(dsc DocumentServiceItem) error {

	qm.AddDocumentService(dsc)
	return nil
}
