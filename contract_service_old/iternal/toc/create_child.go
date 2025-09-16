package toc

import (
	"fmt"
	"regexp"
)

//todo нужно удалить fmt и теперь логируем

// utf16Len возвращает количество UTF-16 code units в строке s

func getChild(re *regexp.Regexp, fullName, parentText, head string, parent TOCItem) []TOCItem {
	var list []TOCItem
	// Находим позиции совпадений
	matches := re.FindAllStringIndex(parentText, -1)
	if len(matches) == 0 {
		fmt.Printf("Нет подзаголовков для '%s' в %s\n", head, parent.Name)
		return nil
	}

	parentStartByte := byteIndexAtRune(fullName, parent.StartPos)
	parentEndByte := byteIndexAtRune(fullName, parent.EndPos)
	parentTextBytes := fullName[parentStartByte:parentEndByte]

	for k, v := range matches {
		var article TOCItem

		relStartRune := runeIndexAtByte(parentTextBytes, v[0])
		relEndRune := runeIndexAtByte(parentTextBytes, v[1])
		absStartRune := parent.StartPos + relStartRune
		absEndRune := 0
		if k < len(matches)-1 {
			nextRelStartRune := runeIndexAtByte(parentTextBytes, matches[k+1][0])
			absEndRune = parent.StartPos + nextRelStartRune - 1
		} else {
			absEndRune = parent.EndPos
		}
		nameStartByte := byteIndexAtRune(fullName, absStartRune)
		nameEndByte := byteIndexAtRune(fullName, parent.StartPos+relEndRune)
		article.Name = fullName[nameStartByte:nameEndByte]
		article.StartPos = absStartRune
		article.EndPos = absEndRune

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

// byteIndexAtRune возвращает байтовый индекс в строке s для rune-позиции pos
func byteIndexAtRune(s string, pos int) int {
	if pos <= 0 {
		return 0
	}
	count := 0
	for i := range s {
		if count == pos {
			return i
		}
		count++
	}
	return len(s)
}

// runeIndexAtByte возвращает rune-индекс в строке s для байтового индекса b
func runeIndexAtByte(s string, b int) int {
	if b <= 0 {
		return 0
	}
	count := 0
	for i := range s {
		if i >= b {
			return count
		}
		count++
	}
	return count
}
