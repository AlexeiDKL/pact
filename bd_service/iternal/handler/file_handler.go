package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"dkl.ru/pact/bd_service/iternal/basedate"
	"dkl.ru/pact/bd_service/iternal/logger"
	"dkl.ru/pact/bd_service/iternal/queue"
)

// структура с зависимостью
type FileHandler struct {
	DB *basedate.Database
	QM *queue.QueueManager
}

// конструктор
func NewFileHandler(db *basedate.Database, qm *queue.QueueManager) *FileHandler {
	return &FileHandler{
		DB: db,
		QM: qm,
	}
}

func (h *FileHandler) SaveFileInBd(w http.ResponseWriter, r *http.Request) {
	h.SaveFile(w, r)
	var req basedate.File
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Невалидный JSON", http.StatusBadRequest)
		return
	}
	if basedate.FileTypeContract == -1 {
		basedate.InitializeFileTypes(h.DB)
	}
	switch req.FileTypeId {
	case basedate.FileTypeContract:
		// добавляем в очередь для создания новой версии, к ней добавляем id только что созданного файла
		return
	case basedate.FileTypeAttachment, basedate.FileTypeOther, basedate.FileTypeFullText, basedate.FileTypeContents:
		// добавляем в очередь, чтобы проверить, нужнно ли его добавить файл в связь с его версией
		return
	default:
		http.Error(w, "Неверный тип файла", http.StatusBadRequest)
		return
	}
}

func (h *FileHandler) SaveFile(w http.ResponseWriter, r *http.Request) {
	var req basedate.File

	// Декодируем тело запроса
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Невалидный JSON", http.StatusBadRequest)
		return
	}

	logger.Logger.Debug(fmt.Sprintf("Получен запрос на сохранение файла: %+v", req))

	err := h.DB.SaveFile(req)
	if err != nil {
		http.Error(w, "Ошибка при сохранении файла", http.StatusInternalServerError)
		logger.Logger.Error(fmt.Sprintf("Ошибка при сохранении файла: %v", err))
		return
	}

	w.WriteHeader(http.StatusCreated)
	logger.Logger.Info(fmt.Sprintf("Файл успешно сохранён: %+v", req))
	w.Write([]byte("✅ Файл успешно сохранён"))
}
