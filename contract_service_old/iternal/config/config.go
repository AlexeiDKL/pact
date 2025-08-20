package config

import (
	"fmt"

	"github.com/spf13/viper"
)

const (
	configName = "config"
	configType = "yaml"
)

var Config ConfigStruct

type ConfigStruct struct {
	Logger LoggerConfig `yaml:"logger"`
	Tokens TokensConfig `yaml:"tokens"`
	BD     BDConfig     `yaml:"bd"`
}

type LoggerConfig struct {
	Path  string `yaml:"path"`
	Type  string `yaml:"type"`
	Level string `yaml:"level"`
	Name  string `yaml:"name"`
}

type TokensConfig struct {
	Garant string `yaml:"garant"`
}

type BDConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}

func Init() error {
	// получаем данные из конфигурационного файла
	// и загружаем их в конфигурацию
	// используем package viper

	viper.SetConfigName(configName)
	viper.SetConfigType(configType)
	viper.AddConfigPath("./contract_service/config/")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	err := viper.Unmarshal(&Config)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}
