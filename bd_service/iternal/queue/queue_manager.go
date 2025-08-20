package queue

func (qm *QueueManager) AddValidation(item ValidationItem) {
	qm.MU.Lock()
	defer qm.MU.Unlock()
	qm.Validation = append(qm.Validation, item)
	qm.ValidationCh <- item // отправляем в канал для обработки
}

func (qm *QueueManager) AddDownload(item DownloadItem) {
	qm.MU.Lock()
	defer qm.MU.Unlock()
	qm.Download = append(qm.Download, item)
	qm.DownloadCh <- item // отправляем в канал для обработки
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

func (qm *QueueManager) RemoveValidationItem(target ValidationItem) {
	qm.MU.Lock()
	defer qm.MU.Unlock()

	newQueue := make([]ValidationItem, 0, len(qm.Validation))
	for _, item := range qm.Validation {
		if item.Body.Topic == target.Body.Topic &&
			item.Body.LanguageID == target.Body.LanguageID &&
			item.Body.FileTypeID == target.Body.FileTypeID {
			continue
		}
		newQueue = append(newQueue, item)
	}
	qm.Validation = newQueue
}

func (qm *QueueManager) RemoveDownloadItem(target DownloadItem) {
	qm.MU.Lock()
	defer qm.MU.Unlock()

	newQueue := make([]DownloadItem, 0, len(qm.Download))
	for _, item := range qm.Download {
		if item.Body.Topic == target.Body.Topic &&
			item.Body.LanguageID == target.Body.LanguageID &&
			item.Body.FileTypeID == target.Body.FileTypeID {
			continue // Пропускаем совпадающий
		}
		newQueue = append(newQueue, item)
	}
	qm.Download = newQueue
}

func (qm *QueueManager) ClearQueues() {
	qm.MU.Lock()
	defer qm.MU.Unlock()
	qm.Validation = nil
	qm.Download = nil
}

func (qm *QueueManager) AddVersion(item BDFile) {
	qm.MU.Lock()
	defer qm.MU.Unlock()
	qm.Version = append(qm.Version, item)
	if qm.VersionCh == nil {
		qm.VersionCh = make(chan BDFile, 100) // инициализируем канал, если он еще не создан
	}
	qm.VersionCh <- item // отправляем в канал для обработки
}

func (qm *QueueManager) RemoveVersionItem(target BDFile) {
	qm.MU.Lock()
	defer qm.MU.Unlock()

	newQueue := make([]BDFile, 0, len(qm.Version))
	for _, item := range qm.Version {
		if item.ID == target.ID && item.LanguageID == target.LanguageID && item.FileTypeID == target.FileTypeID {
			continue // Пропускаем совпадающий
		}
		newQueue = append(newQueue, item)
	}
	qm.Version = newQueue
}

func NewQueueManager() *QueueManager {
	return &QueueManager{
		Validation:   make([]ValidationItem, 0),
		Download:     make([]DownloadItem, 0),
		Version:      make([]BDFile, 0),
		VersionCh:    make(chan BDFile, 100),         // буферизованный канал для версий
		DownloadCh:   make(chan DownloadItem, 100),   // буферизованный канал для скачивания
		ValidationCh: make(chan ValidationItem, 100), // буферизованный канал для валидации
	}
}
