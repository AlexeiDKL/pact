package garant

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"dkl.ru/pact/garant_service/iternal/config"
)

type ModCheckRequest struct {
	Topics     []int64 `json:"topics"`
	ModDate    string  `json:"modDate"` // формат: "YYYY-MM-DD"
	NeedEvents bool    `json:"needEvents"`
}

type ModCheckResponse struct {
	Topics []struct {
		Topic     int64 `json:"topic"`
		ModStatus int   `json:"modStatus"` // 1 - изменился, 2 - не найден
		Events    []struct {
			Date string `json:"date"`
			Type int    `json:"type"` // тип события
		} `json:"events"`
	} `json:"topics"`
}

func CheckModified(topics []int64, sinceDate string) (*ModCheckResponse, error) {
	reqBody := ModCheckRequest{
		Topics:     topics,
		ModDate:    sinceDate,
		NeedEvents: true,
	}

	body, _ := json.Marshal(reqBody)
	req, err := http.NewRequest("POST", "https://api.garant.ru/v1/find-modified", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+config.Config.Tokens.Garant)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ошибка проверки изменений: %s", resp.Status)
	}

	var result ModCheckResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}
