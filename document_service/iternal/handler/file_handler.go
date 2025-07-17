package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"dkl.ru/pact/document_service/iternal/config"
	filesjob "dkl.ru/pact/document_service/iternal/files_job"
)

type FileHandler struct{}

func NewFileHandler() *FileHandler {
	return &FileHandler{}
}

func (h *FileHandler) SaveFileInBd(w http.ResponseWriter, r *http.Request) {}

type GetTextsRequest struct {
	FileName string `json:"file_name"`
	FileID   int    `json:"file_id"`
}

func (h *FileHandler) GetTexts(w http.ResponseWriter, r *http.Request) {
	// получаем имя файла из запроса
	// или id файла из базы данных
	// фозвращаем в ответ текст файла

	var req GetTextsRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Невалидный JSON", http.StatusBadRequest)
		return
	}

	if req.FileName == "" && req.FileID == 0 {
		http.Error(w, "не указано имя файла или ID", http.StatusBadRequest)
		return
	}

	if req.FileName != "" {
		err := h.getTextsByFileName(req.FileName, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	if req.FileID != 0 {
		err := h.getTextsByFileID(req.FileID, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (h *FileHandler) getTextsByFileID(fileID int, w http.ResponseWriter) error {
	fileName := "example.txt" // отправляем запрос в bd_service для получения имени файла по ID
	url := fmt.Sprintf("http://%s:%d/garant/get_file_name_by_id", config.Config.Servers.Bd_service.Host,
		config.Config.Servers.Bd_service.Port)

	payload := map[string]int{"file_id": fileID}

	body, _ := json.Marshal(payload)

	req, err := http.NewRequest("Get", url, bytes.NewReader(body))

	if err != nil {
		return fmt.Errorf("ошибка создания запроса: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return fmt.Errorf("ошибка отправки запроса: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ошибка получения имени файла: статус %d", resp.StatusCode)
	}

	var response map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("ошибка декодирования ответа: %w", err)
	}
	fileName, ok := response["file_name"]
	if !ok {
		return fmt.Errorf("имя файла не найдено в ответе")
	}
	if fileName == "" {
		return fmt.Errorf("имя файла пустое")
	}

	txt, err := filesjob.ReadFileToString(fileName)
	if err != nil {
		return err
	}

	responses := map[string]string{"file_name": fileName, "text": txt}
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(responses)
}

func (h *FileHandler) getTextsByFileName(fileName string, w http.ResponseWriter) error {
	txt, err := filesjob.ReadFileToString(fileName)
	if err != nil {
		return err
	}
	response := map[string]string{"file_name": fileName, "text": txt}
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}
