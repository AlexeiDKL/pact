package httpgarant

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"dkl.ru/pact/contract_service_old/iternal/config"
)

type RequestBody struct {
	Topics  []int  `json:"topics"`
	ModDate string `json:"modDate"`
}

type Response struct {
	Topics []struct {
		Topic     int `json:"topic"`
		ModStatus int `json:"modStatus"`
	} `json:"topics"`
}

// todo покрыть тестами
func CheckFileUpdate(topicID int, modDate string) (bool, error) {
	token := config.Config.Tokens.Garant
	apiURL := "https://api.garant.ru/v1/find-modified"
	// Формируем JSON-тело запроса
	reqBody := RequestBody{
		Topics:  []int{topicID},
		ModDate: modDate,
	}
	body, err := json.Marshal(reqBody)
	if err != nil {
		return false, fmt.Errorf("ошибка кодирования JSON: %v", err)
	}

	// Создаем HTTP-запрос
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(body))
	if err != nil {
		return false, fmt.Errorf("ошибка создания запроса: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	// Отправляем запрос
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("ошибка отправки запроса: %v", err)
	}
	defer resp.Body.Close()

	// Обрабатываем ответ
	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("ошибка API: %v", resp.Status)
	}

	var response Response
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return false, fmt.Errorf("ошибка декодирования ответа: %v", err)
	}

	// Если массив `topics` пустой, возвращаем false
	if len(response.Topics) == 0 {
		return false, nil
	}

	// Проверяем, изменился ли документ
	for _, topic := range response.Topics {
		if topic.Topic == topicID && topic.ModStatus == 1 {
			return true, nil // Документ изменился
		}
	}

	return false, nil // Документ не изменился
}
