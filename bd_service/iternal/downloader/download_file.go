package downloader

import (
	"fmt"
	"io"
	"net/http"
	"os"

	auxiliaryfunctions "dkl.ru/pact/bd_service/iternal/auxiliary_functions"
	"dkl.ru/pact/bd_service/iternal/logger"
)

func downloadFile(url, token, filename string) error {
	logger.Logger.Debug(fmt.Sprintf("try donload with url: %s", url))
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

	// Создание файла
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Копирование содержимого в файл
	_, err = io.Copy(file, resp.Body)
	checksum := auxiliaryfunctions.GetChecksum(file)

	logger.Logger.Debug(fmt.Sprintf("file: %s, downloaded successfully", filename))
	saveFileToDB(checksum)
	return err
}

func saveFileToDB(checksum []byte) {}
