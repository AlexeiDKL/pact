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
	"dkl.ru/pact/bd_service/iternal/handler"
	"dkl.ru/pact/bd_service/iternal/initialization"
	"dkl.ru/pact/bd_service/iternal/logger"
)

func main() {
	db, err := initialization.Init()
	if err != nil {
		panic(err)
	} else {
		logger.Logger.Info("Инициализация успешна")
		logger.Logger.Debug(fmt.Sprintf("%+v", config.Config))
	}

	// 🧩 Инициализация хендлеров
	fileHandler := handler.NewFileHandler(db)
	// 🌐 Создание роутера
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// 📁 Роуты для файлов
	r.Route("/file", func(r chi.Router) {
		r.Post("/", fileHandler.SaveFile) // ✅ Сохранение файла
		// r.Get("/list", fileHandler.GetFilesByVersion)    // ✅ Список файлов по версии
		// r.Get("/meta", fileHandler.GetFileMetaByVersion) // ✅ Дата обновления + ID по версии
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
