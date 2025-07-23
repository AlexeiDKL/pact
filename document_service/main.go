package main

import (
	"log"
	"net/http"

	"dkl.ru/pact/document_service/iternal/handler"
	"dkl.ru/pact/document_service/iternal/initialization"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

/*
📌 Получение текстов документов.
📌 Конвертация файлов из ODT в TXT.
📌 Получение названий "Приложений".
📌 Создание Оглавления.
📌 Создание Полного текста.
*/

func main() {
	// Инициализация конфигурации
	if err := initialization.Init(); err != nil {
		log.Fatalf("Ошибка инициализации: %v", err)
	}

	// Создание роутера
	r := chi.NewRouter()
	r.Use(middleware.Logger) // Используем middleware для логирования запросов

	// Регистрация хендлеров
	fileHandler := handler.NewFileHandler()
	r.Route("/file", func(r chi.Router) {
		r.Post("/get_texts", fileHandler.GetTexts)                 // Получение текстов документов
		r.Post("/convert_odt_to_txt", fileHandler.ConvertOdtToTxt) // Конвертация ODT в TXT
	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("✅ OK, Garant Service is running!"))
	})

	// Запуск HTTP сервера
	http.ListenAndServe(":8082", r)
}
