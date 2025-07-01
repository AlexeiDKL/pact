package handler

import (
	"encoding/json"
	"net/http"

	"dkl.ru/pact/bd_service/iternal/basedate"
)

// структура с зависимостью
type FileHandler struct {
	DB *basedate.Database
}

type SaveFileRequest struct {
	VersionID int    `json:"version_id"`
	FileType  string `json:"file_type"`
	FilePath  string `json:"file_path"`
	Checksum  string `json:"checksum"`
}

// конструктор
func NewFileHandler(db *basedate.Database) *FileHandler {
	return &FileHandler{DB: db}
}

func (h *FileHandler) SaveFile(w http.ResponseWriter, r *http.Request) {
	var req SaveFileRequest

	// Декодируем тело запроса
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Невалидный JSON", http.StatusBadRequest)
		return
	}

	// Вызываем метод сохранения файла
	err := h.DB.SaveFile(req.VersionID, req.FileType, req.FilePath, req.Checksum)
	if err != nil {
		http.Error(w, "Ошибка при сохранении файла", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("✅ Файл успешно сохранён"))
}
