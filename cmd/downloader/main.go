package main

import (
	"log/slog"

	"dkl.ru/pact/cmd/downloader/initiation"
)

/*
	ходим в гарант
	проверяем наличие обновленного текста договора
	Принеобходимости качаем обновленыые файлы
	Записывавем в бд путь к файлам и номер версии

	Для тестов(пока нет интеграции с гарантом) делаем кнопочку, по которой в бд будут закачиваться файлы с маркером новой версии
*/

func initiationPackage() (*slog.Logger, *initiation.Config) {
	config := initiation.InitConfig()
	logger := initiation.InitLogger(config.Log_path, slog.LevelDebug)
	/*
		бд
		API
		сервер
	*/
	return logger, config
}

func downloadFiles() {}

func setVersionInBD() {}
func downloadNewVersion() {
	downloadFiles()
	setVersionInBD()
}

func checkVersion() {}

func main() {
	initiationPackage()
	checkVersion()
	downloadNewVersion()
}
