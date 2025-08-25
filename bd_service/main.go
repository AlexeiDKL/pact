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

	"dkl.ru/pact/bd_service/iternal/basedate"
	"dkl.ru/pact/bd_service/iternal/config"
	garantclient "dkl.ru/pact/bd_service/iternal/garant_client"
	"dkl.ru/pact/bd_service/iternal/handler"
	"dkl.ru/pact/bd_service/iternal/initialization"
	"dkl.ru/pact/bd_service/iternal/logger"
	"dkl.ru/pact/bd_service/iternal/queue"
)

func main() {
	db, err := initialization.Init()
	qm := queue.NewQueueManager()
	// todo дополняем воркеры "корректное закрытие"+ сохранение очереди в файл и заполнение очереди из него
	garantclient.StartDownloadWorker(qm)   // Запускаем воркер для скачивания файлов
	garantclient.StartValidationWorker(qm) // Запускаем воркер для валидации файлов

	basedate.StartVersionWorker(qm, db) // Запускаем воркер для сохранения версий

	// todo воркер, который сканирует бд на наличие version без приложений
	// todo воркер, который сканирует бд на наличие version без полного текста
	// todo воркер, который сканирует бд на наличие version без содержания

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
	// todo заменить на свой middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/file", func(r chi.Router) {
		r.Post("/check_updates", fileHandler.CheckUpdates) // ✅ Проверка обновлений
		r.Get("/download_file", fileHandler.DownloadFile)  // ✅ Скачивание файла todo переименовать
		r.Post("/save", fileHandler.SaveFile)              // ✅ Сохранение файла
		// r.Post("/check", fileHandler.CheckFile)       // ✅ Проверка файла на существование
		// r.Post("/download", fileHandler.DownloadFile) // ✅ Скачивание файла
		// r.Get("/list", fileHandler.GetFilesByVersion)    // ✅ Список файлов по версии
		// r.Get("/meta", fileHandler.GetFileMetaByVersion) // ✅ Дата обновления + ID по версии
	})

	r.Route("/topic", func(r chi.Router) {
		r.Post("/get_language_topics", topicHandler.UpdateTopicsWorkflow) // todo rename url
		r.Post("/set_file_in_bd", fileHandler.SaveFileInBd)               // ✅ Сохранение файла в БД
	})

	// 🔍 Проверка живости сервиса
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("✅ OK"))
	})

	// 🚀 Старт сервера
	logger.Logger.Info(fmt.Sprintf("Сервер запущен на %s:%d\n", config.Config.Server.BdService.Host, config.Config.Server.BdService.Port))
	if err := http.ListenAndServe(
		fmt.Sprintf("%s:%d", config.Config.Server.BdService.Host, config.Config.Server.BdService.Port), r); err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}
}
