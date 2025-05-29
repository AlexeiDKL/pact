package files

import (
	"fmt"
	"net/http"

	"dkl.ru/pact/contract_service/iternal/config"
)

// todo доделать и покрыть тестами
func DownloadFromGarantODT(fileId string) error {
	name := "document"
	fileType := "odt"
	filename := fmt.Sprintf("%s.%s", name, fileType)
	url := fmt.Sprintf("https://api.garant.ru/v1/topic/%s/download-odt", fileId)

	token := config.Config.Tokens.Garant

	return downloadFile(url, token, filename)
}

func downloadFile(url, token, filename string) error {

	return fmt.Errorf("работаем девочки") // todo удаляем строку

	// Создание запроса
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	// Установка заголовков
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

	return Save(filename, resp.Body)
}
