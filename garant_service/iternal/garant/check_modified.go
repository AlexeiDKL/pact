package garant

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"dkl.ru/pact/garant_service/iternal/config"
	"dkl.ru/pact/garant_service/iternal/logger"
)

type ModCheckRequest struct {
	Topics     []int  `json:"topics"`
	ModDate    string `json:"modDate"` // формат: "YYYY-MM-DD"
	NeedEvents bool   `json:"needEvents"`
}

type ModCheckResponse struct {
	Topics []struct {
		Topic     int `json:"topic"`
		ModStatus int `json:"modStatus"` // 1 - изменился, 2 - не найден
		Events    []struct {
			Date string `json:"date"`
			Type int    `json:"type"` // тип события
		} `json:"events"`
	} `json:"topics"`
}

func CheckModified(topics []int, sinceDate string) (*ModCheckResponse, error) {
	reqBody := ModCheckRequest{
		Topics:     topics,
		ModDate:    sinceDate,
		NeedEvents: true,
	}

	logger.Logger.Debug(fmt.Sprintf("Проверяем изменения для тем: %v с даты: %s", topics, sinceDate))

	body, err := json.Marshal(reqBody)
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("Ошибка маршалинга тела запроса: %v", err))
		return nil, err
	}
	logger.Logger.Info(fmt.Sprintf("Отправляем запрос на find-modified: %s", string(body)))

	req, err := http.NewRequest("POST", "https://api.garant.ru/v1/find-modified", bytes.NewReader(body))
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("Ошибка создания запроса: %v", err))
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+config.Config.Tokens.Garant)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("Ошибка выполнения запроса: %v", err))
		return nil, err
	}
	defer resp.Body.Close()

	logger.Logger.Info(fmt.Sprintf("Ответ find-modified: HTTP %d %s", resp.StatusCode, resp.Status))

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("Ошибка чтения тела ответа: %v", err))
		return nil, err
	}
	logger.Logger.Info(fmt.Sprintf("Тело ответа: %s", string(respBody)))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ошибка проверки изменений: %v", resp.Status)
	}

	var result ModCheckResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		logger.Logger.Error(fmt.Sprintf("Ошибка декодирования ответа: %v", err))
		return nil, err
	}
	return &result, nil
}
