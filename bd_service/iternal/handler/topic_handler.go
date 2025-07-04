package handler

import (
	"net/http"

	"dkl.ru/pact/bd_service/iternal/basedate"
	"dkl.ru/pact/bd_service/iternal/core"
	"dkl.ru/pact/bd_service/iternal/logger"
	"dkl.ru/pact/bd_service/iternal/queue"
)

type TopicHandler struct {
	DB *basedate.Database
	QM *queue.QueueManager
}

func NewTopicHandler(db *basedate.Database, qm *queue.QueueManager) *TopicHandler {
	return &TopicHandler{
		DB: db,
		QM: qm,
	}
}

func (h *TopicHandler) GetLanguagesTopics(w http.ResponseWriter, r *http.Request) {
	// Получаем список языков и их топиков из базы данных
	languages, err := h.DB.GetAllLanguages()
	if err != nil {
		http.Error(w, "Ошибка получения языков", http.StatusInternalServerError)
		return
	}

	logger.Logger.Debug("Получен список языков: " + core.MapKeysToString(core.LanguagesToMap(languages)))

	// Инициализируем карту для хранения топиков
	topics := make(map[string]string)
	for _, lang := range languages {
		topic, err := core.GetBaseTopic(*lang.ShortName)
		if err != nil {
			logger.Logger.Error("Ошибка получения топика для языка " + *lang.ShortName + ": " + err.Error())

			continue // Пропускаем язык, если не удалось получить топик
		}
		topics[*lang.ShortName] = topic
	}

	// По  полученным токенам ищем versions и файлы
	logger.Logger.Info("Топики успешно получены для языков: " + core.MapKeysToString(topics))
	if len(topics) == 0 {
		http.Error(w, "Нет доступных топиков", http.StatusNotFound)
		return
	}

	// Извлекаем IDs языков для передачи в GetLatestVersionsByLanguages
	var languageIDs []int
	for _, lang := range languages {
		languageIDs = append(languageIDs, lang.ID)
	}

	res, err := h.DB.GetLatestVersionsByLanguages(languageIDs) // Получаем последние версии по языкам
	if err != nil {
		logger.Logger.Error("Ошибка получения последних версий по языкам: " + err.Error())
		http.Error(w, "Ошибка получения последних версий по языкам", http.StatusInternalServerError)
		return
	}
	logger.Logger.Debug("Получены последние версии по языкам: " + core.MapKeysToString(core.VersionsToMap(res)))

	// Сравниваем, если усть язык, но нет версии, то  топик этого языка добавляем в список на скачивание
	// Если версия есть, топик добавляем в список на проверку актуальности
	//

	for _, lang := range languages {
		found := false
		for _, v := range res {
			if v.LanguageId == lang.ID {
				found = true
				break
			}
		}
		if !found {
			topic, err := core.GetBaseTopic(*lang.ShortName)
			if err != nil {
				logger.Logger.Error("Ошибка получения топика для языка " + *lang.ShortName + ": " + err.Error())
				continue // Пропускаем язык, если не удалось получить топик
			}
			h.QM.AddDownload(queue.DownloadItem{
				LanguageID: lang.ID,
				Topic:      topic,
			})

			logger.Logger.Info("Топик для языка " + *lang.ShortName + " добавлен в список на скачивание: " + topic)
		} else {
			h.QM.AddValidation(queue.ValidationItem{
				LanguageID: lang.ID,
				Topic:      topics[*lang.ShortName],
				VersionID:  res[0].Version, // Используем первую версию, т.к. она последняя
			})
			logger.Logger.Info("Топик для языка " + *lang.ShortName + " добавлен в список на проверку актуальности: " + topics[*lang.ShortName])
		}
	}
	w.WriteHeader(http.StatusOK)
}
