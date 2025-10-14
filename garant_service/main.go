package main

import (
	"fmt"
	"log"
	"net/http"

	"dkl.ru/pact/garant_service/iternal/config"
	documentclient "dkl.ru/pact/garant_service/iternal/document_client"
	"dkl.ru/pact/garant_service/iternal/garant"
	handlers "dkl.ru/pact/garant_service/iternal/handler"
	"dkl.ru/pact/garant_service/iternal/initialization"
	"dkl.ru/pact/garant_service/iternal/logger"
	"dkl.ru/pact/garant_service/iternal/queue"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

/*
📌 Интеграция с API Гаранта.
📌 Поиск документа по названию.
📌 Проверка обновлений текста файлов.
📌 Скачивание файлов из Гаранта.
📌 Механизм повторных попыток скачивания (если API недоступен).
📌 Обработка ошибок API, чтобы избежать "битых" файлов.
*/

func main() {
	err := initialization.Init()
	//todo воркеры для просмотра списков на скачивание и валидацию
	if err != nil {
		panic(err)
	}

	qm := queue.NewQueueManager()

	documentclient.StartConverterWorker(qm)
	queue.StartDownloadWorker(qm)
	queue.StartSaveBDFile(qm)
	queue.StartValidationWorker(qm)

	logger.Logger.Info("Инициализация успешна")
	logger.Logger.Debug("Конфигурация: " + config.Config.String())

	downloadHandler := handlers.NewDownloadListHandler(qm)
	checkHandler := handlers.NewCheckListHandler(qm)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Route("/garant", func(r chi.Router) {
		r.Get("/download", func(w http.ResponseWriter, r *http.Request) {
			err := garant.DownloadODT("70670880", "doc.odt")
			if err != nil {
				http.Error(w, "Ошибка скачивания файла: "+err.Error(), http.StatusInternalServerError)
				return
			}
			w.Write([]byte("✅ Файл успешно скачан и сохранен как doc.odt"))
		})
		r.Post("/add_download", downloadHandler.AddDownloadItem)
		r.Get("/download_list", downloadHandler.GetDownloadList)
		r.Get("/clear_download_list", downloadHandler.ClearDownloadList)

		r.Post("/add_check", checkHandler.AddCheckItem)
		r.Get("/check_list", checkHandler.GetCheckList)
		r.Get("/clear_check_list", checkHandler.ClearCheckList)
	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("✅ OK, Garant Service is running!"))
	})

	logger.Logger.Info(fmt.Sprintf("Сервер запущен на %s:%d\n", config.Config.Server.Garant.Host, config.Config.Server.Garant.Port))
	if err := http.ListenAndServe(
		fmt.Sprintf("%s:%d", config.Config.Server.Garant.Host, config.Config.Server.Garant.Port), r); err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}

}
