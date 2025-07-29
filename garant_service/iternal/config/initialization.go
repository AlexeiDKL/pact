package config

import (
	"fmt"

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
	Server     Servers      `mapstructure:"server"`
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

func (s Servers) String() string {
	return "Servers{" +
		" Garant: " + s.Garant.String() +
		", BdService: " + s.BdService.String() +
		"}"
}

func (c ConfigStruct) String() string {
	return "ConfigStruct{" +
		"Log_config: " + c.Log_config.String() +
		", Tokens: " + c.Tokens.String() +
		", Server: " + c.Server.String() +
		"}"
}

func (s ServerStruct) String() string {
	return fmt.Sprintf("ServerStruct{Host: %s, Port: %d", s.Host, s.Port)
}

type LogStruct struct {
	Path     string `mapstructure:"path_log"`
	LogLevel string `mapstructure:"level_log"`
	Name     string `mapstructure:"name_log"`
	Type     string `mapstructure:"type_log"`
}

func (l LogStruct) String() string {
	return "LogStruct{" +
		"Path: " + l.Path +
		", LogLevel: " + l.LogLevel +
		", Name: " + l.Name +
		", Type: " + l.Type +
		"}"
}

type TokensStruct struct {
	Garant string `mapstructure:"garant"`
}

func (t TokensStruct) String() string {
	return "TokensStruct{" +
		"Garant: " + t.Garant +
		"}"
}

func Init() error {
	// получаем данные из конфигурационного файла
	// и загружаем их в конфигурацию
	// используем package viper

	viper.SetConfigName(configName)
	viper.SetConfigType(configType)
	viper.AddConfigPath("./config/")

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
