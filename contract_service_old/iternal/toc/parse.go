package toc

import (
	"fmt"
	"strings"
)

//todo нужно удалить fmt и теперь логируем

func ParseDocument(fullName, level string, parent TOCItem) []TOCItem {
	var list []TOCItem
	parentText := fullName[parent.StartPos:parent.EndPos]
	fmt.Println(parent)
	if level == "" {
		//pact
		var pact TOCItem
		pact.Name = Header[0]
		pact.StartPos = 0
		pact.EndPos = strings.LastIndex(fullName, "Приложения") - 1
		pact.Children = ParseDocument(fullName, Header[0], pact)

		//appendices
		var appendices TOCItem
		appendices.Name = Header[1]
		appendices.StartPos = strings.LastIndex(fullName, "Приложения")
		appendices.EndPos = len(fullName)
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

	case Header[2]: //Часть
		re := RegexHeader[Header[4]]

		list = append(list, getChild(re, fullName, parentText, Header[4], parent)...)

		// Находим позиции совпадений

		/*
			ищем подчиненные "Раздел"(Header[4]) в parentText
			Объеденяим их в массив list
			Поодному отправляем в ParseDocument(fullName, Header[4], list[i])
		*/

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
