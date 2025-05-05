package initiation

import (
	"log"
	"log/slog"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

// todo перенести в другой покет
type Config struct {
	Env      string
	Log_path string
}

func InitLogger(filename string, level slog.Level) *slog.Logger {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	handler := slog.NewJSONHandler(file, &slog.HandlerOptions{Level: level})

	return slog.New(handler)
}

func InitConfig() *Config {
	os.Setenv("CONFIG_PATH", "./config/local.yaml")
	return MustLoad()
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatalf("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file %s not found", configPath)
	}

	var config Config
	if err := cleanenv.ReadConfig(configPath, &config); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	return &config
}
