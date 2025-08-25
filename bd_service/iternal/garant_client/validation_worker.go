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

func StartValidationWorker(qm *queue.QueueManager) {
	ch := qm.ValidationCh // Получаем канал для валидации
	if ch == nil {
		logger.Logger.Error("❌ Канал для валидации не инициализирован")
		return
	}
	go func() {
		for item := range ch {
			payload := map[string]any{
				"Body": map[string]any{
					"topic":        item.Body.Topic,
					"language_id":  item.Body.LanguageID,
					"version_id":   item.Body.VersionID,
					"file_type_id": item.Body.FileTypeID,
				},
			}

			body, _ := json.Marshal(payload)

			host := config.Config.Server.Garant.Host
			if host == "" {
				host = "localhost"
			}
			url := fmt.Sprintf("http://%s:%d/garant/add_check", host, config.Config.Server.Garant.Port)

			logger.Logger.Debug(fmt.Sprintf("Отправляем запрос на валидацию: %s", url))
			logger.Logger.Debug(fmt.Sprintf("Тело запроса: %s", string(body)))

			req, err := http.NewRequest("POST", url, bytes.NewReader(body))
			if err != nil {
				logger.Logger.Error("❌ [check] Ошибка создания запроса: " + err.Error())
				continue
			}
			req.Header.Set("Content-Type", "application/json")

			resp, err := http.DefaultClient.Do(req)
			if err != nil || resp.StatusCode != http.StatusOK {
				status := "ошибка подключения"
				if err == nil {
					status = resp.Status
				}
				logger.Logger.Error(fmt.Sprintf("⚠️ [check] Ошибка отправки topic %s: %s\n %s", item.Body.Topic, status, url))
				continue
			}

			logger.Logger.Info(fmt.Sprintf("✅ [check] Успешно отправлен файл: %+v", item.Body))
			qm.RemoveValidationItem(item)
		}
	}()
}
