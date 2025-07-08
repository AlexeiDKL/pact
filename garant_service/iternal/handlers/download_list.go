package handlers

import (
	"net/http"

	"dkl.ru/pact/garant_service/iternal/queue"
)

type DownloadListHandler struct {
	QM *queue.QueueManager
}

func NewDownloadListHandler(qm *queue.QueueManager) *DownloadListHandler {
	return &DownloadListHandler{
		QM: qm,
	}
}

func (h *DownloadListHandler) AddDownloadItem(w http.ResponseWriter, r *http.Request) {
	topic := r.URL.Query().Get("topic")
	languageID := r.URL.Query().Get("language_id")
	reason := r.URL.Query().Get("reason")

	item := queue.DownloadItem{
		Topic:      topic,
		LanguageID: languageID,
		Reason:     reason,
	}
	h.QM.AddDownload(item)

	w.Write([]byte("✅ Элемент скачивания добавлен в очередь"))
}

func (h *DownloadListHandler) GetDownloadList(w http.ResponseWriter, r *http.Request) {
	downloadList := h.QM.GetDownloadQueue()
	if len(downloadList) == 0 {
		w.Write([]byte("Список скачивания пуст"))
		return
	}

	for _, item := range downloadList {
		w.Write([]byte(item.Topic + " " + item.LanguageID + " " + item.Reason + "\n"))
	}
}

func (h *DownloadListHandler) ClearDownloadList(w http.ResponseWriter, r *http.Request) {
	h.QM.MU.Lock()
	defer h.QM.MU.Unlock()
	h.QM.Download = make([]queue.DownloadItem, 0)
	w.Write([]byte("✅ Список скачивания очищен"))
}
