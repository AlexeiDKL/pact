package downloader

import (
	"fmt"

	"dkl.ru/pact/bd_service/iternal/config"
	"dkl.ru/pact/bd_service/iternal/logger"
)

func DownloadFromGarantODT(topicId, filename string) {
	url := fmt.Sprintf("https://api.garant.ru/v1/topic/%s/download-odt", topicId)
	token := config.Config.Tokens.Garant

	err := downloadFile(url, token, filename)
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("Ошибка  загрузки файла с id: %s\n ошибка: %s", topicId, err.Error()))
	} else {
		logger.Logger.Debug(fmt.Sprintf("Файл успешно загружен: %s", filename))
	}
}
