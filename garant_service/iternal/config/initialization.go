package config

import (
	myerrors "dkl.ru/pact/garant_service/iternal/my_errors"
	"github.com/spf13/viper"
)

const (
	configName = "config"
	configType = "yaml"
)

var Config ConfigStruct

type ConfigStruct struct {
	Log_config LogStruct    `mapstructure:"log_config"`
	Tokens     TokensStruct `mapstructure:"tokens"`
	Server     ServerStruct `mapstructure:"server"`
}

type ServerStruct struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

type LogStruct struct {
	Path     string `mapstructure:"path_log"`
	LogLevel string `mapstructure:"level_log"`
	Name     string `mapstructure:"name_log"`
	Type     string `mapstructure:"type_log"`
}

type TokensStruct struct {
	Garant string `mapstructure:"garant"`
}

func Init() error {
	// получаем данные из конфигурационного файла
	// и загружаем их в конфигурацию
	// используем package viper

	viper.SetConfigName(configName)
	viper.SetConfigType(configType)
	viper.AddConfigPath("./garant_service/config/")

	if err := viper.ReadInConfig(); err != nil {
		err = myerrors.NotReadConfig(err)
		return err
	}

	err := viper.Unmarshal(&Config)

	if err != nil {
		err = myerrors.NotReadConfig(err)
		return err
	}
	return nil
}
