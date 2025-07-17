package myerrors

type MyError struct {
	Message string
	Err     error
}

func (e *MyError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

func NotReadConfig(err error) error {
	return &MyError{
		Message: "Не удалось прочитать конфигурационный файл",
		Err:     err,
	}
}
