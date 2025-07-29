package config

import (
	"fmt"

	auxiliaryfunctions "dkl.ru/pact/bd_service/iternal/auxiliary_functions"
	myerrors "dkl.ru/pact/bd_service/iternal/my_errors"
	"github.com/spf13/viper"
)

const (
	configName = "config"
	configType = "yaml"
)

var Config ConfigStruct

type ConfigStruct struct {
	Log_config     LogStruct         `mapstructure:"log_config"`
	Bd_server      BDServer          `mapstructure:"bd_server"`
	Tokens         TokensStruct      `mapstructure:"tokens"`
	Document_Topic string            `mapstructure:"document_topic"`
	Server         Servers           `mapstructure:"server"`
	Language_Topic map[string]string `mapstructure:"language_document_topic"`
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

type TokensStruct struct {
	Garant string `mapstructure:"garant"`
}

func (t TokensStruct) String() string {
	return fmt.Sprintf("TokenStruct:{ Garant: %s}",
		auxiliaryfunctions.StringToStar(t.Garant))
}

type LogStruct struct {
	Path     string `mapstructure:"path_log"`
	LogLevel string `mapstructure:"level_log"`
	Name     string `mapstructure:"name_log"`
	Type     string `mapstructure:"type_log"`
}

type BDServer struct {
	Host       string `mapstructure:"host"`
	Port       int    `mapstructure:"port"`
	User       string `mapstructure:"user"`
	Password   string `mapstructure:"password"`
	Dbname     string `mapstructure:"dbname"`
	Path       string `mapstructure:"sql_path"`
	DeletePath string `mapstructure:"sql_delete_path"`
}

func (b BDServer) String() string {
	return fmt.Sprintf("BDServer{Host: %s, Port: %s, User: %s, Password: %s, DBName: %s, Sql_Path: %s, Sql_Delete_Path: %s}",
		auxiliaryfunctions.StringToStar(b.Host),
		auxiliaryfunctions.IntToStar(b.Port), b.User,
		auxiliaryfunctions.StringToStar(b.Password), b.Dbname, b.Path, b.DeletePath)
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
