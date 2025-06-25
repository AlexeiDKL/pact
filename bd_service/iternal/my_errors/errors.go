package myerrors

import "fmt"

func newErrore(message string) error {
	return fmt.Errorf("error: %s", message)
}

func NotRealizeable(name string) error {
	return newErrore(fmt.Sprintf("функция %s не реализована", name))
}

func NotReadConfig(err error) error {
	return newErrore(fmt.Sprintf("Error reading config: %s", err))
}
