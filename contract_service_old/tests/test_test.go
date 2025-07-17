package tests

import (
	"testing"

	config "dkl.ru/pact/contract_service_old/iternal/config"
	"github.com/spf13/viper"
)

func TestInitSuccess(t *testing.T) {
	configName := "config_example"
	configType := "yaml"
	configPath := "./"

	viper.SetConfigName(configName)
	viper.SetConfigType(configType)
	viper.AddConfigPath(configPath)
	viper.AddConfigPath("../config/")

	err := viper.ReadInConfig()
	if err != nil {
		t.Fatalf("Ошибка загрузки конфигурации: %s", err)
	}

	err = viper.Unmarshal(&config.Config)
	if err != nil {
		t.Fatalf("Ошибка парсинга конфигурации: %s", err)
	}
}

// TestInitFileNotFound проверяет обработку ошибки при отсутствии файла
func TestInitFileNotFound(t *testing.T) {
	viper.SetConfigName("nonexistent_config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./")

	err := viper.ReadInConfig()
	if err == nil {
		t.Fatalf("Ожидалась ошибка, но файл загрузился")
	}
}
