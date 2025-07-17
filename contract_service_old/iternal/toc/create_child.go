package toc

import (
	"fmt"
	"regexp"
)

//todo нужно удалить fmt и теперь логируем

func getChild(re *regexp.Regexp, fullName, parentText, head string, parent TOCItem) []TOCItem {
	var list []TOCItem
	// Находим позиции совпадений
	matches := re.FindAllStringSubmatchIndex(parentText, -1)
	if len(matches) == 0 {
		fmt.Printf("Нет подзаголовков для '%s' в %s\n", head, parent.Name)
		return nil
	}

	for k, v := range matches {
		var article TOCItem
		s := parent.StartPos + v[0]
		e := 0
		if k < len(matches)-1 {
			e = parent.StartPos + matches[k+1][0] - 1
		} else {
			e = parent.EndPos
		}
		article.Name = fullName[parent.StartPos+v[0] : parent.StartPos+v[2]-1]
		article.StartPos = s
		article.EndPos = e
		article.Caption = fullName[parent.StartPos+v[2] : parent.StartPos+v[3]]

		switch head {
		case Header[0]:
			article.Children = ParseDocument(fullName, head, article)
		case Header[1]:
			article.Children = ParseDocument(fullName, head, article)
		case Header[2]:
			article.Children = ParseDocument(fullName, head, article)
		case Header[3]:
		case Header[4]:
			article.Children = ParseDocument(fullName, head, article)
		case Header[5]:
		}
		list = append(list, article)

	}
	return list
}
