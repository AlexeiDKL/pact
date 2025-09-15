package router

import (
	"net/http"

	"dkl.ru/pact/orchestrator_service/iternal/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Initialization() *chi.Mux {
	handler := handler.NewSchedulerHandler()
	return setupRouter(handler)
}

func setupRouter(schedulerHandler *handler.SchedulerHandler) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger) // todo заменить на свой логгер

	r.Route("/orchestrator", func(r chi.Router) {
		r.Get("/status", schedulerHandler.Status)
		r.Get("/get_all_versions", schedulerHandler.GetAllVersions)
		r.Post("/get_languages_update_status", schedulerHandler.GetLanguagesUpdateStatus)
		r.Post("/get_file_list", schedulerHandler.GetFileList)
		r.Get("/download_file/{file}", schedulerHandler.DownloadFile)
		r.Post("/create_test_version", schedulerHandler.CreateTestVersion)
		r.Post("/delete_test_version", schedulerHandler.DeleteTestVersion)
	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("✅ OK, Garant Service is running!"))
	})

	return r
}
