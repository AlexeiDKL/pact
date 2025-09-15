package core

import "time"

func CreateNewVersion() int {
	// Получаем текущее время
	now := time.Now()

	// Формируем время начала дня (00:00:00)
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// Получаем Unix timestamp
	timestamp := startOfDay.Unix()
	return int(timestamp)
}
