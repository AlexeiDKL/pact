package handlers

import (
	"encoding/json"
	"net/http"

	"dkl.ru/pact/garant_service/iternal/logger"
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

func (h *CheckListHandler) AddCheckItem(w http.ResponseWriter, r *http.Request) {

	var item queue.ValidationItem
	logger.Logger.Info("Получен запрос на добавление элемента в очередь валидации")
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		logger.Logger.Error("Ошибка декодирования JSON: " + err.Error())
		http.Error(w, "Некорректный JSON", http.StatusBadRequest)
		return
	}
	h.QM.AddValidation(item)

	w.Write([]byte("✅ Элемент валидации добавлен в очередь"))
}

func (h *CheckListHandler) ClearCheckList(w http.ResponseWriter, r *http.Request) {
	h.QM.MU.Lock()
	defer h.QM.MU.Unlock()
	h.QM.Validation = make([]queue.ValidationItem, 0)
	w.Write([]byte("✅ Список валидации очищен"))
}
