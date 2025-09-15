package main

import (
	"fmt"
	"log"
	"net/http"

	"dkl.ru/pact/orchestrator_service/iternal/appinit"
	"dkl.ru/pact/orchestrator_service/iternal/config"
	"dkl.ru/pact/orchestrator_service/iternal/logger"
	"dkl.ru/pact/orchestrator_service/iternal/router"
	"dkl.ru/pact/orchestrator_service/iternal/scheduler"
)

// это сервис, который следит за ограничениями и таймерами,
// вовремя запускает функции в других сервисах.
// А так же переодически проверяет все сервисы на их "здоровье"

func main() {
	// 1. Инициализация приложения
	if err := appinit.Init(); err != nil {
		panic(err)
	}

	// 2. Запуск фоновых задач
	scheduler.StartOrchestrationTasks()

	// 3. Инициализация роутера и хендлеров
	r := router.Initialization()

	// 4. Запуск сервера
	addr := fmt.Sprintf("%s:%d", config.Config.Server.OrchestratorService.Host, config.Config.Server.OrchestratorService.Port)
	logger.Logger.Info(fmt.Sprintf("Сервер запущен на %s", addr))
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}
}
