package initialization

import (
	"dkl.ru/pact/contract_service_old/iternal/config"
	"dkl.ru/pact/contract_service_old/iternal/logger"
)

// todo бд нужно иницилизировать
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
