package main

import (
	"fmt"
	"net/http"

	"dkl.ru/pact/mobile_service/iternal/config"
	"dkl.ru/pact/mobile_service/iternal/handler"
	"dkl.ru/pact/mobile_service/iternal/initialization"
	"dkl.ru/pact/mobile_service/iternal/logger"
	"github.com/go-chi/chi"
)

/*
📌 Проверка обновлений.
📌 Скачивание файлов.
*/

func main() {
	err := initialization.Init()
	if err != nil {
		panic(err)
	}
	logger.Logger.Info("Инициализация успешна")
	logger.Logger.Debug(fmt.Sprintf("Конфигурация: %v", config.Config))

	mobileHandler := handler.NewMobileHandler()

	r := chi.NewRouter()
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("✅ OK, Mobile Service is running!"))
	})
	r.Route("/mobile", func(r chi.Router) {
		r.Get("/check_updates", mobileHandler.CheckUpdates)
		r.Get("/download_file", mobileHandler.DownloadFile)
	})
}
