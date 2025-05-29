package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"dkl.ru/pact/contract_service/iternal/config"
	"dkl.ru/pact/contract_service/iternal/files"
)

var Logger *slog.Logger
var minLogLevel slog.Level

// Кастомный обработчик
type CustomHandler struct {
	file *os.File
}

// Метод Enabled для фильтрации уровней логирования
func (h *CustomHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= minLogLevel
}

// Метод Handle для записи логов в файл
func (h *CustomHandler) Handle(ctx context.Context, r slog.Record) error {
	timestamp := time.Now().Format("2006-01-02 15:04:05") // Формат времени
	level := r.Level.String()                             // Уровень логирования
	msg := r.Message                                      // Сообщение ошибки

	logLine := fmt.Sprintf("%s %s %s\n", timestamp, level, msg)
	_, err := h.file.WriteString(logLine)
	return err
}

// Метод WithAttrs для поддержки атрибутов
func (h *CustomHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h // В данном случае атрибуты не используются
}

// Метод WithGroup для поддержки группировки логов
func (h *CustomHandler) WithGroup(name string) slog.Handler {
	return h // Группировка логов пока не реализована
}

func SetLogLevel(lvl string) {
	switch strings.ToLower(lvl) {
	case "debug":
		minLogLevel = slog.LevelDebug
	case "info":
		minLogLevel = slog.LevelInfo
	case "warn":
		minLogLevel = slog.LevelWarn
	case "error":
		minLogLevel = slog.LevelError
	default:
		minLogLevel = slog.LevelInfo
	}
}

// Функция Init инициализирует логгер и записывает логи в файл
func Init() error {
	logFile := fmt.Sprintf("%s%s.%s", config.Config.Logger.Path, config.Config.Logger.Name, config.Config.Logger.Type)

	if !files.Exists(config.Config.Logger.Path) {
		err := os.MkdirAll(config.Config.Logger.Path, os.ModePerm)
		if err != nil {
			return err
		}
	}

	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {

		return err
	}

	SetLogLevel(config.Config.Logger.Level)

	handler := &CustomHandler{file: file}
	Logger = slog.New(handler) // Инициализация логгера
	Logger.Info("Логгер успешно инициализирован")
	return nil
}
