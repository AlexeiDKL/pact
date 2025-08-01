package documentclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"dkl.ru/pact/garant_service/iternal/config"
	"dkl.ru/pact/garant_service/iternal/core"
	"dkl.ru/pact/garant_service/iternal/files"
	"dkl.ru/pact/garant_service/iternal/logger"
	"dkl.ru/pact/garant_service/iternal/queue"
)

type Response struct {
	Path  string `json:"path"`
	Error string `json:"error"`
}

func StartConverterWorker(qm *queue.QueueManager) {
	ch := qm.DocumentServiceCh
	if ch == nil {
		logger.Logger.Error("❌ Канал для общения с сервисом document не инициализирован")
		return
	}

	go func() {
		/*
			формируем json к document_service:
				путь к файлу
				расширение до
				расширение после
			Пополученному пути формируем запрос к Bd_service, для записи нового файла в бд
		*/
		for item := range ch {

			res, err := handleItem(item)
			if err != nil {
				logger.Logger.Error(err.Error())
				continue
			}

			// Добавляем этот файл в воркер для сохранения в бд
			qm.AddSaveBdFile(res)
			qm.RemoveDocumentServiceItem(item)
		}
	}()
}

func handleItem(item queue.DocumentServiceItem) (queue.BDFile, error) {
	var res queue.BDFile
	path := item.Body.FilePath
	typeBefore := "odt"
	typeAfter := "txt"
	payload := map[string]any{
		"path":   path,
		"before": typeBefore,
		"after":  typeAfter,
	}
	fmt.Println(payload)
	body, _ := json.Marshal(payload)
	host := config.Config.Server.DocumentService.Host
	if host == "" {
		host = "localhost"
	}
	url := fmt.Sprintf("http://%s:%d/file/convert_odt_to_txt", host, config.Config.Server.DocumentService.Port)
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return res, fmt.Errorf("❌ Ошибка создания запроса: " + err.Error())
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return res, fmt.Errorf("⚠️ Ошибка конвертации файла %s: %v", item.Body.FilePath, err)
	}
	defer resp.Body.Close()
	var result Response
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return res, fmt.Errorf("❌ Ошибка при декодировании ответа: " + err.Error())
	}

	// Используем полученные данные
	logger.Logger.Info("✅ Путь к сконвертированному файлу: " + result.Path)
	if result.Error != "" {
		return res, fmt.Errorf("⚠️ Ошибка от сервиса document_service: " + result.Error)
	}

	checksum, err := files.CreateChecksum(result.Path)
	if err != nil {
		return res, fmt.Errorf("⚠️ Ошибка получения checksum для файла %s: %s", res.FilePath, err)
	}

	timestamp, err := core.TimeToStringTimestamp(time.Now())
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("⚠️ Ошибка получения timestamp: %v", err))
	}

	// отправляем файл в bd server, для сохранения в бд и конектимся к version
	res = queue.BDFile{
		ID:           item.Body.ID,
		Checksum:     checksum,
		Name:         result.Path,
		FilePath:     result.Path,
		Topic:        item.Body.Topic,
		LanguageID:   item.Body.LanguageID,
		VersionID:    item.Body.VersionID,
		FileTypeID:   item.Body.FileTypeID,
		DownloadTime: timestamp,
		CreatedAt:    timestamp,
		UpdateAt:     "",
	}

	return res, nil

}
