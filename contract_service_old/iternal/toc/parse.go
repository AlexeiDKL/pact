package toc

import (
	"fmt"
	"regexp"
)

//todo нужно удалить fmt и теперь логируем

func ParseDocument(fullName, level string, parent TOCItem) []TOCItem {
	var list []TOCItem
	parentStartByte := ByteIndexAtRune(fullName, parent.StartPos)
	parentEndByte := ByteIndexAtRune(fullName, parent.EndPos)
	parentText := fullName[parentStartByte:parentEndByte]
	fmt.Println(parent)
	if level == Header[1] { // Приложения
		previewLen := 300
		if len(parentText) < previewLen {
			previewLen = len(parentText)
		}
		fmt.Println("[DEBUG] Превью parentText для 'Приложения':\n" + parentText[:previewLen])
	}
	if level == "" {
		// Найти позицию первого "Приложение № 1[\r\n ]+к Договору о ЕАЭС" (с любым переводом строки)

		markerRe := regexp.MustCompile(`(?i)Приложение\s*[№N]?\s*\d+`)
		loc := markerRe.FindStringIndex(fullName)
		pos := -1
		if loc != nil {
			pos = loc[0]
		} else {
			fmt.Println("[ERROR] Не найден маркер начала приложений (регулярка)")
			// Для отладки:
			previewLen := 200
			if len(fullName) > previewLen {
				fmt.Println("[DEBUG] Начало файла:", fullName[:previewLen])
			}
			pos = len(fullName)
		}

		//pact
		var pact TOCItem
		pact.Name = Header[0]
		pact.StartPos = 0
		pact.EndPos = RuneLen(fullName[:pos-1])
		pact.Children = ParseDocument(fullName, Header[0], pact)

		//appendices
		var appendices TOCItem
		appendices.Name = Header[1]
		appendices.StartPos = RuneLen(fullName[:pos])
		appendices.EndPos = RuneLen(fullName)
		appendices.Children = ParseDocument(fullName, Header[1], appendices)

		list = append(list, pact, appendices)
	}

	switch level {
	case Header[0]: //Договор
		// parentText - текст договора, в котором ищем "Части"
		// Регулярное выражение: ищем строки, где первое слово "Часть", а после него еще одно слово
		re := RegexHeader[Header[2]]

		list = append(list, getChild(re, fullName, parentText, Header[2], parent)...)

	case Header[1]: //Приложения
		// parentText - текст договора, в котором ищем "Части"
		// Регулярное выражение: ищем строки, где первое слово "Часть", а после него еще одно слово

		/*
			ищем подчиненные "Приложение"(Header[3]) в parentText
			Объеденяим их в массив list
			Поодному отправляем в ParseDocument(fullName, Header[3], list[i])
		*/
		re := RegexHeader[Header[3]]

		list = append(list, getChild(re, fullName, parentText, Header[3], parent)...)

	case Header[2]: //Часть
		re := RegexHeader[Header[4]]
		// Находим позиции совпадений

		/*
			ищем подчиненные "Раздел"(Header[4]) в parentText
			Объеденяим их в массив list
			Поодному отправляем в ParseDocument(fullName, Header[4], list[i])
		*/

		list = append(list, getChild(re, fullName, parentText, Header[4], parent)...)

	case Header[3]: //Приложение
		// children nil
		fmt.Println("nil")

	case Header[4]: //Раздел

		re := RegexHeader[Header[5]]

		list = append(list, getChild(re, fullName, parentText, Header[5], parent)...)
		/*
			ищем подчиненные "Статья"(Header[5]) в parentText
			Объеденяим их в массив list
			Поодному отправляем в ParseDocument(fullName, Header[5], list[i])
		*/

	case Header[5]: //Статья

		// children nil
		fmt.Println("nil")
	}
	return list
}
