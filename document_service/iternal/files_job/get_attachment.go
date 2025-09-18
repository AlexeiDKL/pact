package filesjob // todo заголовки получили надо посмотреть, что ещё нужно сделаьть из "контроллеров" и собирать их в единую систему

import (
	"regexp"
	"strings"
)

func GetAttachmentsNames(pactText string) []string {
	var result []string

	// Ищем начала блоков "Приложение ..."
	re := regexp.MustCompile(`(?m)^Приложение\s+(N|№)?\s*(\d+)[\.\s]`)
	locs := re.FindAllStringSubmatchIndex(pactText, -1)
	if len(locs) == 0 {
		return nil
	}
	// Фиктивный конец
	locs = append(locs, []int{len(pactText), len(pactText), len(pactText), len(pactText)})

	// Удаляем только служебные скобки; содержательные оставляем
	serviceBrackets := regexp.MustCompile(`(?i)\s*\((в\s*ред\.?|с\s*изменениями|заголовок\s*в\s*ред\.?|список\s+изменяющих\s+документов).*?\)`)

	// Схлопываем пробелы
	spaceRe := regexp.MustCompile(`\s+`)

	serviceLine := func(s string) bool {
		s = strings.TrimSpace(s)
		return s == "" ||
			strings.HasPrefix(s, "к Договору") ||
			strings.HasPrefix(s, "!!!") ||
			strings.HasPrefix(s, "Список изменяющих документов") ||
			strings.HasPrefix(s, "Список изменяющих документов:")
	}
	clean := func(s string) string {
		s = serviceBrackets.ReplaceAllString(s, "")
		s = strings.TrimSpace(s)
		s = spaceRe.ReplaceAllString(s, " ")
		return s
	}

	for i := 0; i < len(locs)-1; i++ {
		start := locs[i][0]
		end := locs[i+1][0]
		block := pactText[start:end]
		lines := strings.Split(block, "\n")

		// Номер приложения
		numMatch := regexp.MustCompile(`^Приложение\s+(N|№)?\s*(\d+)`).FindStringSubmatch(lines[0])
		if numMatch == nil {
			continue
		}
		num := numMatch[2]

		// Однострочный заголовок в первой строке
		if m := regexp.MustCompile(`^Приложение\s+(N|№)?\s*\d+\.\s*(.+)$`).FindStringSubmatch(lines[0]); m != nil {
			title := clean(m[2]) // только удаляем служебные скобки и нормализуем пробелы
			if title != "" {
				result = append(result, "Приложение № "+num+". "+title)
			}
			continue
		}

		// Многострочный заголовок
		titleLines := []string{}
		foundTitle := false

		for j := 1; j < len(lines); j++ {
			line := strings.TrimSpace(lines[j])

			// Стоп/служебные строки до начала названия — пропускаем
			if !foundTitle {
				if serviceLine(line) {
					continue
				}
				cl := clean(line)
				if cl != "" {
					titleLines = append(titleLines, cl)
					foundTitle = true
				}
				continue
			}

			// Уже начали название: конец на пустой/служебной строке
			if serviceLine(line) {
				break
			}

			// Добавляем следующую строку заголовка, очищая только служебные скобки
			cl := clean(line)
			if cl != "" {
				titleLines = append(titleLines, cl)
			}
		}

		if len(titleLines) > 0 {
			title := clean(strings.Join(titleLines, " "))
			if title != "" {
				result = append(result, "Приложение № "+num+". "+title)
			}
		}
	}

	return result
}
