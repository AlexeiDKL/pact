package handlers

import (
	"net/http"

	"dkl.ru/pact/garant_service/iternal/queue"
)

type CheckListHandler struct {
	QM *queue.QueueManager
}

func NewCheckListHandler(qm *queue.QueueManager) *CheckListHandler {
	return &CheckListHandler{
		QM: qm,
	}
}

func (h *CheckListHandler) AddCheckItem(w http.ResponseWriter, r *http.Request) {

	topic := r.URL.Query().Get("topic")
	languageID := r.URL.Query().Get("language_id")
	versionID := r.URL.Query().Get("version_id")
	fileType := r.URL.Query().Get("file_type")

	item := queue.ValidationItem{
		Topic:      topic,
		LanguageID: languageID,
		VersionID:  versionID,
		FileType:   fileType,
	}
	h.QM.AddValidation(item)

	w.Write([]byte("✅ Элемент валидации добавлен в очередь"))
}

func (h *CheckListHandler) GetCheckList(w http.ResponseWriter, r *http.Request) {
	checkList := h.QM.GetValidationQueue()
	if len(checkList) == 0 {
		w.Write([]byte("Список валидации пуст"))
		return
	}

	for _, item := range checkList {
		w.Write([]byte(item.Topic + " " + item.LanguageID + " " + item.VersionID + " " + item.FileType + "\n"))
	}
}
