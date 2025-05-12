package initiation

import (
	"database/sql"
	"log"
	"log/slog"
	"os"

	mssql "github.com/denisenkom/go-mssqldb"
	"github.com/ilyakaznacheev/cleanenv"
)

// todo перенести в другой покет
type Config struct {
	Env      string
	Log_path string
}

func InitLogger(filename string, level slog.Level) (*slog.Logger, error) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	handler := slog.NewJSONHandler(file, &slog.HandlerOptions{Level: level})

	return slog.New(handler), nil
}

func InitConfig() (*Config, error) {
	os.Setenv("CONFIG_PATH", "./config/local.yaml")
	return MustLoadConfig(), nil
}

func MustLoadConfig() *Config {
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

func InitBd() (*sql.DB, error) {

	_, err := mssql.NewConnector("")
	if err != nil {
		return nil, err
	}
	return nil, nil
}
