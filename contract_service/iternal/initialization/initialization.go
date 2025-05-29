package initialization

import (
	"dkl.ru/pact/contract_service/iternal/config"
	"dkl.ru/pact/contract_service/iternal/logger"
)

func Init() error {
	err := config.Init()
	if err != nil {
		return err
	}
	err = logger.Init()
	if err != nil {
		return err
	}
	return nil
}
