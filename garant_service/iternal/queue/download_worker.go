package queue

import (
	"fmt"

	"dkl.ru/pact/garant_service/iternal/garant"
	"dkl.ru/pact/garant_service/iternal/logger"
)

func StartDownloadWorker(qm *QueueManager) {

	ch := qm.DownloadCh
	if ch == nil {
		logger.Logger.Error("❌ Канал для скачивания не инициализирован")
		return
	}
	go func() {
		/*
			выделяем из qm элемент
			качаем из гаранта файл
			отправляем в document_service
		*/
		for item := range ch {
			topic := item.Topic
			fmt.Println(item)
			// создайм путь
			// определяем язык и версию
			var fileType string
			var version string
			var language string
			if item.LanguageID == "" {
				language = "ru" // по умолчанию русский
			}
			if item.VersionID == "" {
				version = "1" // по умолчанию первая версия
			}
			if item.FileType == "" {
				// по определенному языку пишем тип на нужном языке
				fileType = "договор" // по умолчанию договор
			}
			fmt.Println("Получаем файл для темы:", topic, "язык:", language, "версия:", version, "тип файла:", fileType)
			fileName := fmt.Sprintf("./files/%s_%s.odt", fileType, version)
			err := garant.DownloadODT(topic, fileName)
			if err != nil {
				logger.Logger.Error(fmt.Sprintf("❌ Ошибка скачивания файла для темы %s: %v", topic, err))
				continue // пропускаем этот элемент и продолжаем цикл
			}
			// добавляем в очередь, которую обработает воркер
			// воркер отправляет информацию о файле в document_service
			// для этого добавляем в qm новый элемент
			// qm DocumentService     []DocumentServiceItem,	DocumentServiceCh   chan DocumentServiceItem

			documentServiceItem := DocumentServiceItem{
				Topic:       topic,
				LanguageID:  item.LanguageID,
				FileType:    fileType,
				FileVersion: version,
				FileName:    fileName,
			}
			err = qm.SendFileToDocumentService(documentServiceItem)
			if err != nil {
				logger.Logger.Error(fmt.Sprintf("❌ Ошибка отправки файла для темы %s: %v", topic, err))
				continue // пропускаем этот элемент и продолжаем цикл
			}

			qm.RemoveDownloadItem(item) // удаляем элемент из очереди после успешной отправки
		}
	}()
}
