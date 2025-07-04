package main

/*
	📌 Инициализация готово
		Конфиг
		Лог
		Бд
			Иницилизируем бд.
				Проверяем наличие бд, при необходимости создаем
				Проверяем наличие таблиц, при необходимости создаем
	📌 Получение файлов. готово
	📌 Получение списка файлов. список файлов из бд?
			Получаем список файлов из бд по версии. готово
	📌 Получение даты обновления файла и его ID.
			Получаем дату обновления файла и его ID по версии.
	📌 Сохранение файла в БД. готово
	📌 Получение новых версий.
	📌 Обновление информации о файлах.

*/

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"dkl.ru/pact/bd_service/iternal/config"
	"dkl.ru/pact/bd_service/iternal/core"
	"dkl.ru/pact/bd_service/iternal/handler"
	"dkl.ru/pact/bd_service/iternal/initialization"
	"dkl.ru/pact/bd_service/iternal/logger"
	"dkl.ru/pact/bd_service/iternal/queue"
)

func main() {
	db, err := initialization.Init()
	qm := queue.NewQueueManager()
	if err != nil {
		panic(err)
	} else {
		logger.Logger.Info("Инициализация успешна")
		logger.Logger.Debug(fmt.Sprintf("%+v", config.Config))
	}

	// 🧩 Инициализация хендлеров
	fileHandler := handler.NewFileHandler(db, qm)

	topicHandler := handler.NewTopicHandler(db, qm)

	// 🌐 Создание роутера
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// 📁 Роуты для файлов
	r.Route("/v1", func(r chi.Router) {
		r.Get("/get_queue", func(w http.ResponseWriter, r *http.Request) {
			// Получаем очередь на валидацию и скачивание
			validationQueue := qm.GetValidationQueue()
			downloadQueue := qm.GetDownloadQueue()

			// Формируем ответ
			response := map[string]interface{}{
				"validation": validationQueue,
				"download":   downloadQueue,
			}

			if err := core.WriteJSONResponse(w, response); err != nil {
				logger.Logger.Error("Ошибка при отправке ответа: " + err.Error())
				http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
			}
		})
		r.Post("/clear_queue", func(w http.ResponseWriter, r *http.Request) {
			// Очищаем очереди
			qm.ClearQueues()
			logger.Logger.Info("Очереди успешно очищены")

			// Отправляем ответ
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("✅ Очереди успешно очищены"))
		})
	})
	r.Route("/file", func(r chi.Router) {
		r.Post("/", fileHandler.SaveFile)             // ✅ Сохранение файла
		r.Post("/check", fileHandler.CheckFile)       // ✅ Проверка файла на существование
		r.Post("/download", fileHandler.DownloadFile) // ✅ Скачивание файла
		// r.Get("/list", fileHandler.GetFilesByVersion)    // ✅ Список файлов по версии
		// r.Get("/meta", fileHandler.GetFileMetaByVersion) // ✅ Дата обновления + ID по версии
	})

	r.Route("/topic", func(r chi.Router) {
		r.Post("/get_language_topics", topicHandler.GetLanguagesTopics) // todo rename ✅ Получение файлов по языку
	})

	// 🔧 Заглушки для будущих фич
	// r.Get("/version/new", fileHandler.GetNewVersions)
	// r.Put("/file", fileHandler.UpdateFileInfo)

	// 🔍 Проверка живости сервиса
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("✅ OK"))
	})

	// 🚀 Старт сервера
	log.Println("Сервер запущен на http://localhost:8080")
	if err := http.ListenAndServe(
		fmt.Sprintf("%s:%d", config.Config.Server.Host, config.Config.Server.Port), r); err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}
}
