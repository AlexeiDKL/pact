package garantclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"dkl.ru/pact/bd_service/iternal/config"
	"dkl.ru/pact/bd_service/iternal/logger"
	"dkl.ru/pact/bd_service/iternal/queue"
)

func StartDownloadWorker(qm *queue.QueueManager) {
	ch := qm.DownloadCh // Получаем канал для скачивания
	if ch == nil {
		logger.Logger.Error("❌ Канал для скачивания не инициализирован")
		return
	}
	go func() {
		for item := range ch {
			payload := map[string]any{
				"topic":      item.Body.Topic,
				"LanguageID": item.Body.LanguageID,
				"VersionID":  item.Body.VersionID, // timestamp от 00:00
				"FileType":   item.Body.FileTypeID,
			}

			body, _ := json.Marshal(payload)
			host := config.Config.Server.Garant.Host
			if host == "" {
				host = "localhost"
			}
			url := fmt.Sprintf("http://%s:%d/garant/add_download", host, config.Config.Server.Garant.Port)

			req, err := http.NewRequest("POST", url, bytes.NewReader(body))
			if err != nil {
				logger.Logger.Error("❌ Ошибка создания запроса: " + err.Error())
				continue
			}
			req.Header.Set("Content-Type", "application/json")

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				logger.Logger.Error(fmt.Sprintf("⚠️ Ошибка отправки topic %s: %v", item.Body.Topic, err))
				continue
			}
			if resp.StatusCode != http.StatusOK {
				logger.Logger.Error(fmt.Sprintf("⚠️ Ошибка отправки topic %s  statusCode: %d url: %s", item.Body.Topic, resp.StatusCode, url))
				continue
			}

			logger.Logger.Info("✅ Успешно отправлен topic на скачивание: " + item.Body.Topic)
			qm.RemoveDownloadItem(item)
		}
	}()
}
