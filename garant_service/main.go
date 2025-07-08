package main

import (
	"fmt"
	"net/http"

	"dkl.ru/pact/garant_service/iternal/config"
	"dkl.ru/pact/garant_service/iternal/garant"
	"dkl.ru/pact/garant_service/iternal/initialization"
	"github.com/go-chi/chi/v5"
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
	if err != nil {
		panic(err)
	}

	topic := "70670880"
	err = garant.DownloadODT(topic, "doc.odt")
	if err != nil {
		fmt.Printf("Ошибка скачивания файла: %v\n", err)

		return
	}
	fmt.Println("✅ Файл успешно скачан и сохранен как doc.odt")

	r := chi.NewRouter()

	r.Route("/garant", func(r chi.Router) {
		r.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello, Garant Service!"))
		})
		r.Get("/download", func(w http.ResponseWriter, r *http.Request) {
			err := garant.DownloadODT("70670880", "doc.odt")
			if err != nil {
				http.Error(w, "Ошибка скачивания файла: "+err.Error(), http.StatusInternalServerError)
				return
			}
			w.Write([]byte("✅ Файл успешно скачан и сохранен как doc.odt"))
		})
	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("✅ OK, Garant Service is running!"))
	})

	fmt.Printf("Starting Garant Service on port %s\n", config.Config.Server.Port)
	if err := http.ListenAndServe(fmt.Sprintf("%s:%s", config.Config.Server.Host, config.Config.Server.Port), r); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
		return
	}

}
