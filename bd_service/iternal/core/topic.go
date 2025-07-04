package core

import (
	"strings"

	"dkl.ru/pact/bd_service/iternal/config"
	myerrors "dkl.ru/pact/bd_service/iternal/my_errors"
)

func GetBaseTopic(language string) (string, error) {
	language = strings.ToLower(language)

	topic := strings.ToLower(config.Config.Language_Topic[language])
	if topic == "" {
		return "", myerrors.ErrLanguageNotFound(language)
	}
	return topic, nil
}
