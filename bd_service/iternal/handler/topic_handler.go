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

func (h *TopicHandler) UpdateTopicsWorkflow(w http.ResponseWriter, r *http.Request) {
	// Получаем список языков и их топиков из базы данных
	languages, err := h.DB.GetAllLanguages()
	if err != nil {
		http.Error(w, "Ошибка получения языков", http.StatusInternalServerError)
		return
	}

	logger.Logger.Debug("Получен список языков: " + core.MapKeysToString(core.LanguagesToMap(languages)))

	// Инициализируем мапу для хранения топиков
	topics := make(map[string]string)
	for _, lang := range languages {
		topic, err := core.GetBaseTopicFromConfig(*lang.ShortName)
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

	res, err := h.DB.GetLatestVersionsByLanguagesID(languageIDs) // Получаем последние версии по языкам
	if err != nil {
		logger.Logger.Error("Ошибка получения последних версий по языкам: " + err.Error())
		http.Error(w, "Ошибка получения последних версий по языкам", http.StatusInternalServerError)
		return
	}
	logger.Logger.Debug("Получены последние версии по языкам: " + core.MapKeysToString(core.VersionsToMap(res)))

	// Сравниваем, если eсть язык, но нет версии, то  топик этого языка добавляем в список на скачивание
	// Если версия есть, топик добавляем в список на проверку актуальности
	//

	for _, lang := range languages {
		foundVersion := false
		var versionItem basedate.Version
		for _, v := range res {
			if v.LanguageId == lang.ID {
				foundVersion = true
				versionItem = v
				break
			}
		}
		topic, err := core.GetBaseTopicFromConfig(*lang.ShortName)
		if err != nil {
			logger.Logger.Error("Ошибка получения топика для языка " + *lang.ShortName + ": " + err.Error())
			continue
		}
		ftypeName := "pact"
		fTypeId, err := h.DB.GetFyleTypeByName(ftypeName)
		if err != nil {
			logger.Logger.Error("Ошибка получения id типа файла " + ftypeName + ": " + err.Error())
		}
		if !foundVersion {
			h.QM.AddDownload(queue.DownloadItem{
				Body: queue.BDFile{
					Topic:      topic,
					LanguageID: lang.ID,
					FileTypeID: fTypeId,
				},
			})
			logger.Logger.Info("Топик для языка " + *lang.ShortName + " добавлен в список на скачивание")
		} else {

			// передаём запросом в сервис garant_service
			// В garant_service подумать, какие данные передавать, для одназначного определения файла
			// topic - уникальный id файла, используем для проверки актуальности
			// Данные для идентификации файла -
			// LanguageID - ID языка
			// VersionID - ID версиия к которой он относится
			// FileType - тип файла, наприме "договор", "приложение", "полный текст"
			h.QM.AddValidation(queue.ValidationItem{
				Body: queue.BDFile{
					Topic:      topic,
					LanguageID: lang.ID,
					FileTypeID: fTypeId,
					VersionID:  int(versionItem.Version),
				},
			})
			logger.Logger.Info("Топик для языка " + *lang.ShortName + " добавлен в список на проверку актуальности: " + topics[*lang.ShortName])
		}
	}
	w.WriteHeader(http.StatusOK)
}
