package basedate

import (
	"fmt"

	"dkl.ru/pact/bd_service/iternal/logger"
	"dkl.ru/pact/bd_service/iternal/queue"
)

func StartVersionWorker(qm *queue.QueueManager, db *Database) {
	logger.Logger.Info("Запуск воркера для сохранения версий")
	ch := qm.VersionCh // Получаем канал для версий
	if ch == nil {
		logger.Logger.Error("❌ Канал для версий не инициализирован")
		return
	}
	go func() {
		for item := range ch {
			send := File{
				Id:         item.ID,
				FileTypeId: item.FileTypeID,
				LanguageId: item.LanguageID,
				CreatedAt:  item.CreatedAt,
				UpdatedAt:  item.UpdatedAt,
				Name:       item.Name,
			}
			err := db.SaveToVersion(send)
			if err != nil {
				logger.Logger.Error(fmt.Sprintf("❌ Ошибка сохранения версии: %v", err))
				continue
			}
			logger.Logger.Info(fmt.Sprintf("✅ Версия успешно сохранена: %+v", item))
			qm.RemoveVersionItem(item) // Удаляем элемент из очереди после успешного сохранения
		}
	}()
}
