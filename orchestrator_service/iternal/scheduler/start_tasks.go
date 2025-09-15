package scheduler

import (
	"time"

	"dkl.ru/pact/orchestrator_service/iternal/clock"
	"dkl.ru/pact/orchestrator_service/iternal/config"
	"dkl.ru/pact/orchestrator_service/iternal/logger"
)

func StartOrchestrationTasks() {
	worker := NewScheduledWorker(
		time.Duration(config.Config.Orchestrator.SyncIntervalMs),
		time.Duration(config.Config.Orchestrator.SyncStartTimeMs),
		func() {
			logger.Logger.Debug("Запуск фоновой задачи синхронизации...")
			// todo запуск синхронизации http://localhost:8080/topic/get_language_topics
		},
		clock.NewRealClock(),
	)
	chekAttachmentsWorker := NewScheduledWorker(
		time.Duration(config.Config.Orchestrator.SyncIntervalMs),
		time.Duration(config.Config.Orchestrator.SyncStartTimeMs+(2*60*60*1000)), // на 2 часа позже
		func() {
			logger.Logger.Debug("Запуск фоновой задачи проверки приложений...")
			// todo запуск проверки вложений http://localhost:8080/topic/chek_attachments
		},
		clock.NewRealClock(),
	)
	healthWorker := NewScheduledWorker(
		time.Duration(config.Config.Orchestrator.SyncIntervalMs),
		time.Duration(config.Config.Orchestrator.SyncStartTimeMs/2),
		func() {
			logger.Logger.Debug("Запуск фоновой задачи проверки здоровья сервисов...")
			//todo запуск проверки здоровья сервисов /orchestrator/status

		},
		clock.NewRealClock(),
	)

	worker.Start()
	healthWorker.Start()
	chekAttachmentsWorker.Start()
	defer worker.Stop()
	defer healthWorker.Stop()
	defer chekAttachmentsWorker.Stop()
}
