package config

import (
	myerrors "dkl.ru/pact/document_service/iternal/my_errors"
	"github.com/spf13/viper"
)

const (
	configName = "config"
	configType = "yaml"
)

var Config ConfigStruct

type ConfigStruct struct {
	Log_config LogStruct `mapstructure:"log_config"`
	Servers    Servers   `mapstructure:"server"`
}

type Servers struct {
	Garant_service   ServerStruct `mapstructure:"garant_service"`
	Bd_service       ServerStruct `mapstructure:"bd_service"`
	Document_service ServerStruct `mapstructure:"document_service"`
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

func Init() error {
	// получаем данные из конфигурационного файла
	// и загружаем их в конфигурацию
	// используем package viper

	viper.SetConfigName(configName)
	viper.SetConfigType(configType)
	viper.AddConfigPath("./document_service/config/")

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
