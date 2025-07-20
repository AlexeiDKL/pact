package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"dkl.ru/pact/mobile_service/iternal/config"
	"dkl.ru/pact/mobile_service/iternal/logger"
)

type MobileHandler struct{}

func NewMobileHandler() *MobileHandler {
	return &MobileHandler{}
}

func (h *MobileHandler) CheckUpdates(w http.ResponseWriter, r *http.Request) {
	type updateRequest struct {
		Language string `json:"language"`
		Version  string `json:"version"`
	}
	type updateResponse struct {
		UpdateAvailable bool   `json:"update_available"`
		Error           string `json:"error,omitempty"`
	}

	// 1. Прочитать тело запроса целиком
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, updateResponse{Error: "не удалось прочитать тело запроса"})
		return
	}
	defer r.Body.Close()

	// 2. Декодировать JSON
	var ur updateRequest
	if err := json.Unmarshal(bodyBytes, &ur); err != nil {
		writeJSON(w, http.StatusBadRequest, updateResponse{Error: "неверный формат JSON"})
		return
	}
	if ur.Language == "" || ur.Version == "" {
		writeJSON(w, http.StatusBadRequest, updateResponse{Error: "отсутствует язык или версия"})
		return
	}

	// 3. Подготовить запрос к bd_service (POST с JSON)
	url := fmt.Sprintf(
		"http://%s:%d/file/check_updates",
		config.Config.Server.BdService.Host,
		config.Config.Server.BdService.Port,
	)
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(bodyBytes))
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("ошибка создания запроса к bd_service: %v", err))
		writeJSON(w, http.StatusInternalServerError, updateResponse{Error: "внутренняя ошибка"})
		return
	}
	req.Header.Set("Content-Type", "application/json")

	// 4. Выполнить запрос с таймаутом
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Logger.Error(fmt.Sprintf("запрос к bd_service не выполнен: %v", err))
		writeJSON(w, http.StatusBadGateway, updateResponse{Error: "не удалось связаться с сервисом обновлений"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Logger.Error(fmt.Sprintf("bd_service вернул статус %d", resp.StatusCode))
		writeJSON(w, resp.StatusCode, updateResponse{Error: "ошибка сервиса обновлений"})
		return
	}

	// 5. Декодировать ответ (ожидаем JSON { "update_available": true/false })
	var body map[string]bool
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		logger.Logger.Error(fmt.Sprintf("декодирование ответа bd_service: %v", err))
		writeJSON(w, http.StatusInternalServerError, updateResponse{Error: "не удалось разобрать ответ"})
		return
	}

	// 6. Отдать клиенту финальный ответ
	writeJSON(w, http.StatusOK, updateResponse{UpdateAvailable: body["update_available"]})
}

// Утилита для единообразной отправки JSON-ответов
func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func (h *MobileHandler) DownloadFile(w http.ResponseWriter, r *http.Request) {
	// Получаем язык
	// Полуаем в bd_service послееднию версию для этого языка, а так же пути к полному тексту и оглавлению
	// Отпраляем эти файлы на клиент
	// todo
	language := r.URL.Query().Get("language")
	if language == "" {
		http.Error(w, "Отсутствует язык", http.StatusBadRequest)
		return
	}
}
