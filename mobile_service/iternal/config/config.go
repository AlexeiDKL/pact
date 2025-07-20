package config

import (
	myerrors "dkl.ru/pact/mobile_service/iternal/my_errors"
	"github.com/spf13/viper"
)

const (
	configName = "config"
	configType = "yaml"
)

var Config ConfigStruct

type ConfigStruct struct {
	Log_config LogStruct `mapstructure:"log_config"`
	Server     Servers   `mapstructure:"server"`
}

type Servers struct {
	Garant          ServerStruct `mapstructure:"garant_service"`
	BdService       ServerStruct `mapstructure:"bd_service"`
	MobileService   ServerStruct `mapstructure:"mobile_service"`
	DocumentService ServerStruct `mapstructure:"document_service"`
}

type ServerStruct struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
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
	viper.AddConfigPath("./mobile_service/config/")

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
