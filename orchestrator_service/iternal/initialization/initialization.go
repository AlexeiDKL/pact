package initialization

import (
	"fmt"

	"dkl.ru/pact/orchestrator_service/iternal/config"
	"dkl.ru/pact/orchestrator_service/iternal/logger"
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
	logger.Logger.Debug(fmt.Sprintf("%+v", config.Config))
	return nil
}
