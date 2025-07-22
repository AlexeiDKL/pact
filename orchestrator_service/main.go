package main

import (
	"dkl.ru/pact/orchestrator_service/iternal/handler"
	"dkl.ru/pact/orchestrator_service/iternal/initialization"
	"dkl.ru/pact/orchestrator_service/iternal/logger"
	"github.com/go-chi/chi/v5"
)

// это сервис, который следит за ограничениями и таймерами,
// вовремя запускает функции в других сервисах.
// А так же переодически проверяет все сервисы на их "здоровье"

func main() {
	err := initialization.Init()
	if err != nil {
		panic(err)
	}

	schedulerHandler := handler.NewSchedulerHandler()
	logger.Logger.Debug("Ку")

	r := chi.NewRouter()

	r.Route("/orchestrator", func(r chi.Router) {
		r.Post("/start-sync", schedulerHandler.StartSync)
		r.Post("/retry-failed", schedulerHandler.RetryFailed)
		r.Post("/shutdown", schedulerHandler.Shutdown)
		r.Get("/status", schedulerHandler.Status)
	})
}
