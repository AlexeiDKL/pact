package myerrors

import "fmt"

func newErrore(message string) error {
	return fmt.Errorf("error: %s", message)
}

func NotDownload(name string, err error) error {
	return fmt.Errorf("не удалось скачать %s: %w", name, err)
}

func NotRealizeable(name string) error { //todo на момент тестирования во всех сервисах должен возвращать панику, на релизе, только ошибку
	panic(fmt.Sprintf("функция %s не реализована", name))
	// return newErrore(fmt.Sprintf("функция %s не реализована", name))
}

func NotReadConfig(err error) error {
	return newErrore(fmt.Sprintf("Error reading config: %s", err))
}

func ErrLanguageNotFound(language string) error {
	return newErrore(fmt.Sprintf("язык %s не найден в конфигурации", language))
}
