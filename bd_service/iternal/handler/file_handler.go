package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"dkl.ru/pact/bd_service/iternal/basedate"
	"dkl.ru/pact/bd_service/iternal/logger"
)

// структура с зависимостью
type FileHandler struct {
	DB *basedate.Database
}

// конструктор
func NewFileHandler(db *basedate.Database) *FileHandler {
	return &FileHandler{DB: db}
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
