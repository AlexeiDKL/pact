package core

import (
	"fmt"

	"dkl.ru/pact/bd_service/iternal/basedate"
)

func MapKeysToString(m map[string]string) string {
	if len(m) == 0 {
		return "Нет доступных топиков"
	}

	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}

	return "Топики: " + joinKeys(keys)
}

func joinKeys(keys []string) string {
	if len(keys) == 0 {
		return ""
	}
	result := keys[0]
	for _, key := range keys[1:] {
		result += ", " + key
	}
	return result
}

func LanguagesToMap(languages []basedate.Language) map[string]string {
	if len(languages) == 0 {
		return nil
	}

	languagesMap := make(map[string]string, len(languages))
	for _, lang := range languages {
		languagesMap[*lang.ShortName] = lang.FullName
	}

	return languagesMap
}

func VersionsToMap(versions []basedate.Version) map[string]string {
	if len(versions) == 0 {
		return nil
	}
	versionsMap := make(map[string]string, len(versions))
	for _, v := range versions {
		versionsMap[fmt.Sprintf("%d", v.Version)] = fmt.Sprintf("%d", v.PactId)
	}

	return versionsMap
}
