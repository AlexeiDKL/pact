package documentclient

import (
	"dkl.ru/pact/garant_service/iternal/logger"
	"dkl.ru/pact/garant_service/iternal/queue"
)

func StartConverterWorker(qm *queue.QueueManager) { //todo 25.07.25
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
			в document_service:
				читаем содержимое файла
				Если файла нет возвращаем ошибку
				создаём путь к "конечному файлу", удаляем расширение у пути к файлу, и подставляем целевое расширение
				через switch_case выполняем запуск необходимой функции в которую передаём прочитанный файл и путь до "итогового файл"
		*/
	}()
}
