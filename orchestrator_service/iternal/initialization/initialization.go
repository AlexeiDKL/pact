package initialization

import (
	"dkl.ru/pact/orchestrator_service/iternal/config"
	"dkl.ru/pact/orchestrator_service/iternal/logger"
	myerrors "dkl.ru/pact/orchestrator_service/iternal/my_errors"
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
	return myerrors.NotRealizeable("dkl.ru/pact/orchestrator_service/iternal/initialization")
}
