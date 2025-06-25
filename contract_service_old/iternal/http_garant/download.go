package httpgarant

import (
	"fmt"
	"net/http"

	"dkl.ru/pact/contract_service/iternal/config"
	"dkl.ru/pact/contract_service/iternal/files"
)

// todo покрыть тестами
func DownloadFromGarantODT(fileId string) error {
	name := "document"
	fileType := "odt"
	filename := fmt.Sprintf("%s.%s", name, fileType)
	url := fmt.Sprintf("https://api.garant.ru/v1/topic/%s/download-odt", fileId)

	return downloadFile(url, filename)
}

func downloadFile(url, filename string) error {
	token := config.Config.Tokens.Garant

	// Создание запроса
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	// Установка заголовков
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	// Выполнение запроса
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Проверка ответа
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ошибка: статус %d", resp.StatusCode)
	}

	return files.Save(filename, resp.Body)
}
