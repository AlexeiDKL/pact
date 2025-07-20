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
	case basedate.FileTypeContract: // todo
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

// Типы запроса и ответа вынесены на уровень пакета
type CheckUpdatesRequest struct {
	Language string `json:"language"`
	Version  string `json:"version"`
}

type CheckUpdatesResponse struct {
	UpdateAvailable bool   `json:"update_available"`
	Error           string `json:"error,omitempty"`
}

func (h *FileHandler) CheckUpdates(w http.ResponseWriter, r *http.Request) {
	// 1. Разрешён только POST
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, CheckUpdatesResponse{
			Error: "только POST метод",
		})
		return
	}

	// 2. Декодируем тело и закрываем ридер
	var req CheckUpdatesRequest
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, CheckUpdatesResponse{
			Error: "неверный формат JSON",
		})
		return
	}
	if req.Language == "" || req.Version == "" {
		writeJSON(w, http.StatusBadRequest, CheckUpdatesResponse{
			Error: "не указаны язык или версия",
		})
		return
	}

	// 3. Передаём контекст в DB для отмены при таймауте
	ctx := r.Context()
	ok, err := h.DB.CheckUpdates(ctx, req.Language, req.Version)
	if err != nil {
		logger.Logger.Error("CheckUpdates error:", err)
		writeJSON(w, http.StatusInternalServerError, CheckUpdatesResponse{
			Error: "ошибка при проверке в БД",
		})
		return
	}

	// 4. Отдаём единый JSON-ответ
	writeJSON(w, http.StatusOK, CheckUpdatesResponse{
		UpdateAvailable: ok,
	})
	logger.Logger.Info("CheckUpdates:", req.Language, req.Version, ok)
}

// Помощник для отправки JSON
func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
