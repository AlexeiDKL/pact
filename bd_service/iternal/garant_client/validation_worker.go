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
				"topic":      item.Body.Topic,
				"LanguageID": item.Body.LanguageID,
				"VersionID":  item.Body.VersionID,
				"FileType":   item.Body.FileTypeID,
			}

			body, _ := json.Marshal(payload)
			url := fmt.Sprintf("http://%s:%s/garant/add_check", config.Config.Server.Garant.Host, config.Config.Server.Garant.Port)

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
				logger.Logger.Error(fmt.Sprintf("⚠️ [check] Ошибка отправки topic %s: %s", item.Body.Topic, status))
				continue
			}

			logger.Logger.Info("✅ [check] Успешно отправлен topic: " + item.Body.Topic)
			qm.RemoveValidationItem(item)
		}
	}()
}
