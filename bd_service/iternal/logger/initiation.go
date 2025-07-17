package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"dkl.ru/pact/bd_service/iternal/config"
	"dkl.ru/pact/contract_service_old/iternal/files" // todo переделать
)

var Logger *slog.Logger
var minLogLevel slog.Level

//todo логировать в JSON

// Кастомный обработчик
type CustomHandler struct {
	file *os.File
}

// Метод Enabled для фильтрации уровней логирования
func (h *CustomHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= minLogLevel
}

// Структура для упорядочивания JSON-логов
type LogEntry struct {
	Time  string `json:"time"`
	Level string `json:"level"`
	Msg   string `json:"msg"`
	Stack string `json:"stack,omitempty"`
}

// Метод Handle для записи логов в файл
func (h *CustomHandler) Handle(ctx context.Context, r slog.Record) error {
	timestamp := time.Now().Format("2006-01-02 15:04:05") // Формат времени
	var stack string
	if r.Level >= slog.LevelError {
		stack = string(debug.Stack()) // ← получаем стек
	}

	entry := LogEntry{
		Time:  timestamp,
		Level: r.Level.String(),
		Msg:   r.Message,
		Stack: stack,
	}

	jsonData, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("ошибка кодирования JSON: %v", err)
	}

	_, err = h.file.WriteString(string(jsonData) + "\n")
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
	logFile := fmt.Sprintf("%s%s.%s", config.Config.Log_config.Path, config.Config.Log_config.Name, config.Config.Log_config.Type)

	if !files.Exists(config.Config.Log_config.Path) {
		err := os.MkdirAll(config.Config.Log_config.Path, os.ModePerm)
		if err != nil {
			return err
		}
	}

	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {

		return err
	}

	SetLogLevel(config.Config.Log_config.LogLevel)

	handler := &CustomHandler{file: file}
	Logger = slog.New(handler) // Инициализация логгера
	Logger.Info("Логгер успешно инициализирован")
	return nil
}
