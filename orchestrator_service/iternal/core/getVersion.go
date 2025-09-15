package core

import "strings"

func GetVersions(allVerson map[string]map[string]any) map[string]int {
	result := make(map[string]int)
	for lang, data := range allVerson {
		ver, ok := data["version"].(int)
		if !ok {
			// если версия хранится как float64 (например, после json.Unmarshal)
			if vf, okf := data["version"].(float64); okf {
				ver = int(vf)
			} else {
				continue
			}
		}
		result[strings.ToLower(lang)] = ver
	}
	return result
}
