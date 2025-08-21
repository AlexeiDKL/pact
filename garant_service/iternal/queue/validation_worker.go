package queue

import (
	"strconv"
	"time"

	"dkl.ru/pact/garant_service/iternal/garant"
	"dkl.ru/pact/garant_service/iternal/logger"
)

func StartValidationWorker(qm *QueueManager) {
	ch := qm.ValidationCh // Получаем канал для валидации
	if ch == nil {
		logger.Logger.Error("❌ Канал для валидации не инициализирован")
		return
	}
	go func() {
		for item := range ch {
			// Да - Отправляем запрос в Гарант, по `topic` для проверки актуальности
			logger.Logger.Debug("Получен запрос на валидацию: " + item.Body.Topic)
			// Здесь должна быть логика отправки запроса на валидацию
			now := time.Now().Format("2025-04-14")
			i, _ := strconv.Atoi(item.Body.Topic)

			res, err := garant.CheckModified([]int{i}, now)
			if err != nil {
				logger.Logger.Error("Ошибка при проверке актуальности: " + err.Error())
				continue
			}
			if res.Topics[0].ModStatus == 1 {
				logger.Logger.Info("Тема " + item.Body.Topic + " актуальна, отправляем на валидацию")
				// Отправляем на валидацию
				send := DownloadItem(item)
				qm.AddDownload(send)
			} else {
				logger.Logger.Info("Тема " + item.Body.Topic + " не актуальна, пропускаем")
			}
		}
	}()
}
