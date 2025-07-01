package main

/*
📌 Инициализация готово
	Конфиг
	Лог
	Бд
		Иницилизируем бд.
			Проверяем наличие бд, при необходимости создаем
			Проверяем наличие таблиц, при необходимости создаем
📌 Получение файлов. готово
📌 Получение списка файлов. список файлов из бд?
📌 Получение даты обновления файла и его ID.
📌 Сохранение файла в БД.
📌 Получение новых версий.
📌 Обновление информации о файлах.
*/

import (
	"fmt"

	"dkl.ru/pact/bd_service/iternal/config"
	"dkl.ru/pact/bd_service/iternal/downloader"
	"dkl.ru/pact/bd_service/iternal/initialization"
	"dkl.ru/pact/bd_service/iternal/logger"
)

func main() {
	_, err := initialization.Init()
	if err != nil {
		panic(err)
	} else {
		logger.Logger.Info("Инициализация успешна")
		logger.Logger.Debug(fmt.Sprintf("%+v", config.Config))
	}
	checksum, err := downloader.DownloadFromGarantODT(config.Config.Document_Topic, "doc.odt")
	if err != nil {
		panic(err)
	}

	fmt.Println(checksum)

}
