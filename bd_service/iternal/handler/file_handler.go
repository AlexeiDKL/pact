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

func (h *FileHandler) CheckFile(w http.ResponseWriter, r *http.Request) {
	// Получаем из qm очередь на валидацию
	// Проверяем актуальность файла по api Гарант
	// Если файл актуален, то возвращаем 200 OK
	// Если файл не актуален, то добавляем его в список на скачивание
	// и возвращаем 202 Accepted
}

func (h *FileHandler) DownloadFile(w http.ResponseWriter, r *http.Request) {
	// Получаем из qm очередь на скачивание
	// Скачиваем файл по api Гарант
	// Если файл успешно скачан, то записываем его в базу данных
	// Проверяем тип файла, если это "pact", то создаём новую версию, если другой, то добавляем связи в бд
	// Если файл не найден, то возвращаем 404 Not Found
	// Если произошла ошибка при скачивании, то возвращаем 500 Internal Server Error
}
