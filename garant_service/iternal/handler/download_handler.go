package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"dkl.ru/pact/garant_service/iternal/logger"
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
	logger.Logger.Info("Получен запрос на добавление элемента в очередь скачивания")
	logger.Logger.Debug(fmt.Sprintf("Тело запроса: %v", r.Body))
	var item queue.DownloadItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		logger.Logger.Error("Ошибка декодирования JSON: " + err.Error())
		http.Error(w, "Некорректный JSON", http.StatusBadRequest)
		return
	}
	h.QM.AddDownload(item)
	logger.Logger.Info("Добавлен элемент в очередь скачивания: " + item.Topic + " " + item.LanguageID)
	w.Write([]byte("✅ Элемент скачивания добавлен в очередь"))
}

func (h *DownloadListHandler) GetDownloadList(w http.ResponseWriter, r *http.Request) {
	downloadList := h.QM.GetDownloadQueue()
	if len(downloadList) == 0 {
		w.Write([]byte("Список скачивания пуст"))
		return
	}

	for _, item := range downloadList {
		w.Write([]byte(item.Topic + " " + item.LanguageID + " " + item.VersionID + "\n"))
	}
}

func (h *DownloadListHandler) ClearDownloadList(w http.ResponseWriter, r *http.Request) {
	h.QM.MU.Lock()
	defer h.QM.MU.Unlock()
	h.QM.Download = make([]queue.DownloadItem, 0)
	w.Write([]byte("✅ Список скачивания очищен"))
}
