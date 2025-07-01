package initialization

import (
	"fmt"

	"dkl.ru/pact/bd_service/iternal/basedate"
	"dkl.ru/pact/bd_service/iternal/config"
	"dkl.ru/pact/bd_service/iternal/logger"
)

func Init() (*basedate.Database, error) {
	//config
	err := config.Init()
	if err != nil {
		return nil, err
	}

	//log
	if err := logger.Init(); err != nil {
		return nil, err
	} else {
		logger.Logger.Debug("Конфиг инициализирован")
		logger.Logger.Debug("Логгер инициализирован")
	}

	//bd
	db, err := basedate.New(config.Config.Bd_server, logger.Logger)
	if err != nil {
		logger.Logger.Warn(fmt.Sprintf("Ошибка инициализации БД: %v", err))
	} else {
		logger.Logger.Debug("БазаДанных инициализирована")
	}

	return db, err
}
