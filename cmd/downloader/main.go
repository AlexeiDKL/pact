package main

import (
	"database/sql"
	"fmt"
	"log/slog"

	"dkl.ru/pact/cmd/downloader/initiation"
)

/*
	ходим в гарант
	проверяем наличие обновленного текста договора
	Принеобходимости качаем обновленыые файлы
	Записывавем в бд путь к файлам и номер версии

	Для тестов(пока нет интеграции с гарантом) делаем кнопочку, по которой в бд будут закачиваться файлы с маркером новой версии



	Ещё проверочка конфига
*/

func initiationPackage() (*slog.Logger, *initiation.Config, *sql.DB, error) {
	config, err := initiation.InitConfig()
	if err != nil {
		return nil, nil, nil, err
	}
	logger, err := initiation.InitLogger(config.Log_path, slog.LevelDebug)
	if err != nil {
		return nil, nil, nil, err
	}
	bd, err := initiation.InitBd()
	if err != nil {
		fmt.Println(err.Error())
		return nil, nil, nil, err
	}

	/*
		API
		сервер
	*/
	return logger, config, bd, nil
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
