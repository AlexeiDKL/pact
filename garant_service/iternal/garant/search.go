package garant

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"dkl.ru/pact/contract_service_old/iternal/config"
)

func SearchFile(fileName string) (Document, error) {
	url := "https://api.garant.ru/v1/search"
	token := config.Config.Tokens.Garant
	data := Data{
		Text:      fmt.Sprintf("& BOOL(& MorphoName(%s))", fileName),
		Count:     30,
		IsQuery:   true,
		Kind:      []string{"001", "002", "003"},
		Sort:      0,
		SortOrder: 0,
	}
	jsonData, err := json.Marshal(data)

	var result Document
	if err != nil {
		return result, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return result, err
	}

	// return nil
	// Установка заголовков
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	// Выполнение запроса
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	// Проверка ответа
	if resp.StatusCode != http.StatusOK {
		return result, fmt.Errorf("ошибка: статус %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return result, err
	}

	var docs Documents
	err = json.Unmarshal(body, &docs)
	if err != nil {
		return result, err
	}

	result = docs.Documents[0]
	// resp.Body результат запроса
	return result, nil
}
