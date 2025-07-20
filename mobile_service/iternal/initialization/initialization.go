package initialization

import (
	"fmt"

	"dkl.ru/pact/mobile_service/iternal/config"
	"dkl.ru/pact/mobile_service/iternal/logger"
)

func Init() error {
	// Инициализация конфигурации
	if err := config.Init(); err != nil {
		return fmt.Errorf("ошибка инициализации конфигурации: %v", err)
	}

	// Инициализация логгера
	if err := logger.Init(); err != nil {
		return fmt.Errorf("ошибка инициализации логгера: %v", err)
	}

	logger.Logger.Info("Сервис успешно инициализирован")
	return nil
}
