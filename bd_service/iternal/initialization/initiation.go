package initialization

import (
	"dkl.ru/pact/bd_service/iternal/basedate"
	"dkl.ru/pact/bd_service/iternal/config"
	"dkl.ru/pact/bd_service/iternal/logger"
)

func Init() error {
	//config
	err := config.Init()
	if err != nil {
		return err
	}

	//log
	if err := logger.Init(); err != nil {
		return err
	} else {
		logger.Logger.Debug("Конфиг инициализирован")
		logger.Logger.Debug("Логгер инициализирован")
	}

	//bd
	if err := basedate.Init(); err != nil {
		return err
	} else {
		logger.Logger.Debug("БазаДанных инициализирована")
	}
	return nil
}
