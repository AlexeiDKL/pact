package handler

import (
	"encoding/json"
	"fmt"
	"io"
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
       logger.Logger.Info("SaveFileInBd: получен запрос на сохранение файла")
       h.SaveFile(w, r)
       var req basedate.File
       if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
	       logger.Logger.Error("SaveFileInBd: невалидный JSON: ", err)
	       http.Error(w, "Невалидный JSON", http.StatusBadRequest)
	       return
       }
       logger.Logger.Debug(fmt.Sprintf("SaveFileInBd: параметры файла: %+v", req))
       if basedate.FileTypePact == -1 {
	       logger.Logger.Info("SaveFileInBd: инициализация типов файлов")
	       basedate.InitializeFileTypes(h.DB)
       }
       switch req.FileTypeId {
       case basedate.FileTypePact:
	       logger.Logger.Info("SaveFileInBd: файл типа Pact, добавление в очередь для создания новой версии")
	       // добавляем в очередь для создания новой версии, к ней добавляем id только что созданного файла
	       return
       case basedate.FileTypeAttachment, basedate.FileTypeOther, basedate.FileTypeFullText, basedate.FileTypeContents:
	       logger.Logger.Info("SaveFileInBd: файл типа Attachment/Other/FullText/Contents, добавление в очередь для проверки связи с версией")
	       // добавляем в очередь, чтобы проверить, нужнно ли его добавить файл в связь с его версией
	       return
       default:
	       logger.Logger.Warn("SaveFileInBd: неверный тип файла: ", req.FileTypeId)
	       http.Error(w, "Неверный тип файла", http.StatusBadRequest)
	       return
       }
}

func (h *FileHandler) SaveFile(w http.ResponseWriter, r *http.Request) {
       var req basedate.File

       if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
	       logger.Logger.Error("SaveFile: невалидный JSON: ", err)
	       http.Error(w, "Невалидный JSON", http.StatusBadRequest)
	       return
       }

       logger.Logger.Debug(fmt.Sprintf("SaveFile: получен запрос на сохранение файла: %+v", req))

       id, err := h.DB.SaveFile(req)
       if err != nil {
	       logger.Logger.Error(fmt.Sprintf("SaveFile: ошибка при сохранении файла: %v", err))
	       http.Error(w, "Ошибка при сохранении файла", http.StatusInternalServerError)
	       return
       }

       req.Id = id

       item := queue.BDFile{
	       ID:         req.Id,
	       FileTypeID: req.FileTypeId,
	       LanguageID: req.LanguageId,
	       CreatedAt:  req.CreatedAt,
	       UpdatedAt:  req.UpdatedAt,
	       Name:       req.Name,
       }

       logger.Logger.Debug(fmt.Sprintf("SaveFile: добавление файла в очередь: %+v", item))
       h.QM.AddVersion(item)

       w.WriteHeader(http.StatusCreated)
       logger.Logger.Info(fmt.Sprintf("✅ Файл успешно сохранён: %+v", req))
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
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		 logger.Logger.Error("writeJSON: ошибка при кодировании JSON: ", err)
	}
}

func (h *FileHandler) DownloadFile(w http.ResponseWriter, r *http.Request) {
       logger.Logger.Info("DownloadFile: получен запрос на скачивание файла")
       type downloadFileResponse struct {
	       FullFilePath    string `json:"full_file_path"`
	       ContentFilePath string `json:"content_file_path"`
	       Error           string `json:"error,omitempty"`
       }
       type downloadFileRequest struct {
	       Language string `json:"language"`
       }
       bodyBytes, err := io.ReadAll(r.Body)
       if err != nil {
	       logger.Logger.Error("DownloadFile: не удалось прочитать тело запроса: ", err)
	       writeJSON(w, http.StatusBadRequest, downloadFileResponse{Error: "не удалось прочитать тело запроса"})
	       return
       }
       defer r.Body.Close()

       var dr downloadFileRequest
       if err := json.Unmarshal(bodyBytes, &dr); err != nil {
	       logger.Logger.Error("DownloadFile: неверный формат JSON: ", err)
	       writeJSON(w, http.StatusBadRequest, downloadFileResponse{Error: "неверный формат JSON"})
	       return
       }
       if dr.Language == "" {
	       logger.Logger.Warn("DownloadFile: отсутствует язык в запросе")
	       writeJSON(w, http.StatusBadRequest, downloadFileResponse{Error: "отсутствует язык"})
	       return
       }
       ctx := r.Context()
       logger.Logger.Debug(fmt.Sprintf("DownloadFile: получаем пути файлов для языка %s", dr.Language))
       pathFiles, err := h.DB.GetPathFileClasVersionByLanguage(ctx, dr.Language)

       if err != nil {
	       logger.Logger.Error("DownloadFile: ошибка получения пути к файлу: ", err)
	       writeJSON(w, http.StatusInternalServerError, downloadFileResponse{Error: "ошибка получения пути к файлу"})
	       return
       }
       if pathFiles.Full == "" || pathFiles.Content == "" {
	       logger.Logger.Warn(fmt.Sprintf("DownloadFile: нет файлов для языка %s", dr.Language))
	       writeJSON(w, http.StatusNotFound, downloadFileResponse{Error: "нет файлов для языка " + dr.Language})
	       return
       }
       writeJSON(w, http.StatusOK, downloadFileResponse{
	       FullFilePath:    pathFiles.Full,
	       ContentFilePath: pathFiles.Content,
       })
       logger.Logger.Info(fmt.Sprintf("DownloadFile: файлы успешно отправлены для языка %s: полный текст %s, оглавление %s", dr.Language, pathFiles.Full, pathFiles.Content))
}

func (h *FileHandler) GetFileNameByID(w http.ResponseWriter, r *http.Request) {
       logger.Logger.Info("GetFileNameByID: получен запрос на получение имени файла по ID")
       fileID := r.URL.Query().Get("file_id")
       if fileID == "" {
	       logger.Logger.Warn("GetFileNameByID: file_id параметр обязателен")
	       http.Error(w, "file_id параметр обязателен", http.StatusBadRequest)
	       return
       }
       var id int
       _, err := fmt.Sscanf(fileID, "%d", &id)
       if err != nil {
	       logger.Logger.Warn("GetFileNameByID: file_id должен быть числом")
	       http.Error(w, "file_id должен быть числом", http.StatusBadRequest)
	       return
       }
       logger.Logger.Debug(fmt.Sprintf("GetFileNameByID: ищем имя файла для file_id %d", id))
       fileName, err := h.DB.GetFileNameByID(r.Context(), id)
       if err != nil {
	       logger.Logger.Error(fmt.Sprintf("GetFileNameByID: ошибка при получении имени файла: %v", err))
	       http.Error(w, "ошибка при получении имени файла: "+err.Error(), http.StatusInternalServerError)
	       return
       }
       type getFileNameResponse struct {
	       FileName string `json:"file_name"`
       }
       resp := getFileNameResponse{FileName: fileName}
       w.Header().Set("Content-Type", "application/json")
       if err := json.NewEncoder(w).Encode(resp); err != nil {
	       logger.Logger.Error("GetFileNameByID: ошибка при кодировании ответа: ", err)
	       http.Error(w, "ошибка при кодировании ответа: "+err.Error(), http.StatusInternalServerError)
	       return
       }
       logger.Logger.Info(fmt.Sprintf("GetFileNameByID: имя файла успешно получено для file_id %d: %s", id, fileName))
}
