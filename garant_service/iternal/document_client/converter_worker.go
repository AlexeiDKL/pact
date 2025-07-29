package documentclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"dkl.ru/pact/garant_service/iternal/config"
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
			fmt.Printf("%v", item)
			err := handleItem(item)
			if err != nil {
				logger.Logger.Error(err.Error())
				continue
			}

			// Добавляем этот файл в воркер для сохранения в бд
			qm.AddSaveBdFile(item.Body)
			qm.RemoveDocumentServiceItem(item)
		}
	}()
}

func handleItem(item queue.DocumentServiceItem) error {
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
		return fmt.Errorf("❌ Ошибка создания запроса: " + err.Error())
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("⚠️ Ошибка конвертации файла %s: %v", item.Body.FilePath, err)
	}
	defer resp.Body.Close()
	var result Response
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return fmt.Errorf("❌ Ошибка при декодировании ответа: " + err.Error())
	}

	// Используем полученные данные
	logger.Logger.Info("✅ Путь к сконвертированному файлу: " + result.Path)
	if result.Error != "" {
		return fmt.Errorf("⚠️ Ошибка от сервиса document_service: " + result.Error)
	}

	return nil

}
