package initialization

import (
	"dkl.ru/pact/document_service/iternal/config"
	"dkl.ru/pact/document_service/iternal/logger"
)

func Init() error {
	// config Init
	// logger init
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
