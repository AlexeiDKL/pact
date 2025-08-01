package queue

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"dkl.ru/pact/garant_service/iternal/config"
	"dkl.ru/pact/garant_service/iternal/logger"
)

func StartSaveBDFile(qm *QueueManager) {
	ch := qm.SaveBDFileCH
	if ch == nil {
		logger.Logger.Error("❌ Канал для сохранения в бд не инициализирован")
		return
	}
	go func() {
		/*
			отправляем данные в bd service
			для сохранения файла в бд и дальнейших действий с ними
		*/
		for item := range ch {
			body, err := json.Marshal(item)
			if err != nil {
				logger.Logger.Error(fmt.Sprintf("Не удалось конвертировать в json: %s", err.Error()))
				continue
			}
			host := config.Config.Server.BdService.Host
			if host == "" {
				host = "localhost"
			}
			port := config.Config.Server.BdService.Port
			url := fmt.Sprintf("http://%s:%d/file/save", host, port)
			req, err := http.NewRequest("POST", url, bytes.NewReader(body))
			if err != nil {
				logger.Logger.Error("❌ Ошибка создания запроса: " + err.Error())
				continue
			}
			req.Header.Set("Content-Type", "application/json")

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				logger.Logger.Error(fmt.Sprintf("⚠️ Ошибка отправки файла в бд %s: %v", item.Name, err))
				continue
			}
			if resp.StatusCode != http.StatusOK {
				logger.Logger.Error(fmt.Sprintf("⚠️ Ошибка отправки файла в бд %s  statusCode: %d url: %s", item.Name, resp.StatusCode, url))
				continue
			}
			logger.Logger.Info("✅ Успешно отправлен файл в бд: " + item.Name)
			qm.RemoveSaveBdFile(item)
		}
	}()
}
