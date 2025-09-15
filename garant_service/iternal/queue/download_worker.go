package queue

import (
	"fmt"

	"dkl.ru/pact/garant_service/iternal/core"
	"dkl.ru/pact/garant_service/iternal/files"
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
			topic := item.Body.Topic
			fmt.Println(item)
			// создайм путь
			// определяем язык и версию
			fileType := item.Body.FileType
			version := item.Body.VersionID
			language := item.Body.LanguageID
			fmt.Println(item)
			if item.Body.LanguageID == 0 {
				language = 5 // по умолчанию русский
			}
			if item.Body.VersionID == -1 {
				version = core.CreateNewVersion()
			}
			if item.Body.FileType == -1 {
				// по определенному языку пишем тип на нужном языке
				fileType = 1 // по умолчанию договор
			}
			fmt.Println("Получаем файл для темы:", topic, "язык:", language, "версия:", version, "тип файла:", fileType)
			fileName := fmt.Sprintf("../files/%d_%d.odt", fileType, version)

			if !files.FileExists(fileName) {
				err := garant.DownloadODT(topic, fileName)
				if err != nil {
					logger.Logger.Error(fmt.Sprintf("❌ Ошибка скачивания файла для темы %s: %v", topic, err))
					continue // пропускаем этот элемент и продолжаем цикл
				}
			}
			// добавляем в очередь, которую обработает воркер
			// воркер отправляет информацию о файле в document_service
			// для этого добавляем в qm новый элемент
			// qm DocumentService     []DocumentServiceItem,	DocumentServiceCh   chan DocumentServiceItem

			documentServiceItem := DocumentServiceItem{
				Body: BDFile{
					Name:       fileName,
					FilePath:   fileName,
					Topic:      topic,
					VersionID:  version,
					LanguageID: item.Body.LanguageID,
					FileTypeID: fileType,
				},
			}
			err := qm.SendFileToDocumentService(documentServiceItem)
			if err != nil {
				logger.Logger.Error(fmt.Sprintf("❌ Ошибка отправки файла для темы %s: %v", topic, err))
				continue // пропускаем этот элемент и продолжаем цикл
			}

			qm.RemoveDownloadItem(item) // удаляем элемент из очереди после успешной отправки
		}
	}()
}
