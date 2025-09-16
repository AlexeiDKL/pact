package toc

import (
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"
)

// Настрой здесь: номер обязателен
var strictRe = regexp.MustCompile(`(?im)^[\p{Zs}\t]*(приложение[\p{Zs}\t]*(?:№|N)?[\p{Zs}\t]*\d+)[\p{Zs}\t]*[^\r\n]*`)

// Широкая сетка-кандидат (любой "прилож" от начала строки)
var candRe = regexp.MustCompile(`(?im)^[^\S\r\n]*прилож`)

func otladka(raw string) {
	// 0) Нормализуем окончания строк и BOM, чтобы позиции и ^/$ были предсказуемыми
	text := strings.TrimPrefix(raw, "\uFEFF")
	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")

	// A) Кандидаты (очень широкая сетка)
	cands := candRe.FindAllStringIndex(text, -1)

	// B) Принято строгим шаблоном (с подмасками)
	accepted := strictRe.FindAllStringSubmatchIndex(text, -1)

	fmt.Printf("Diag: candidates=%d, accepted=%d\n", len(cands), len(accepted))

	// Печатаем ВСЕ принятые совпадения подробно
	for i, m := range accepted {
		// m: [fullStart, fullEnd, g1Start, g1End, ...]
		fullStart, fullEnd := m[0], m[1]
		var headStart, headEnd int
		if len(m) >= 4 && m[2] >= 0 && m[3] >= 0 {
			headStart, headEnd = m[2], m[3]
		} else {
			// На случай отсутствия группы 1
			headStart, headEnd = fullStart, fullEnd
		}

		lineStart := lastNL(text, fullStart) + 1
		lineEnd := nextNL(text, fullStart)
		if lineEnd < 0 {
			lineEnd = len(text)
		}
		line := text[lineStart:lineEnd]
		lineNo := strings.Count(text[:lineStart], "\n") + 1

		header := strings.TrimSpace(text[headStart:headEnd])

		fmt.Printf("ACCEPT[%02d] line %d @%d..%d | head %d..%d\n", i, lineNo, fullStart, fullEnd, headStart, headEnd)
		fmt.Printf("  header: %q\n", header)
		fmt.Printf("  line:   %q\n", visualize(line))
	}

	// Индекс начальных позиций принятых (по началу полной строки)
	acceptedStarts := make(map[int]struct{}, len(accepted))
	for _, m := range accepted {
		acceptedStarts[m[0]] = struct{}{}
	}

	// C) Разница: кандидаты, которые не прошли строгий regex
	for _, c := range cands {
		candStart := c[0]
		if _, ok := acceptedStarts[candStart]; ok {
			continue
		}
		lineStart := lastNL(text, candStart) + 1
		lineEnd := nextNL(text, candStart)
		if lineEnd < 0 {
			lineEnd = len(text)
		}
		line := text[lineStart:lineEnd]
		lineNo := strings.Count(text[:lineStart], "\n") + 1

		// Следующая строка (часто бывает "к Договору …")
		nextStart := lineEnd + 1
		nextEnd := nextNL(text, nextStart)
		if nextEnd < 0 {
			nextEnd = len(text)
		}
		nextLine := ""
		if nextStart >= 0 && nextStart < len(text) {
			nextLine = text[nextStart:nextEnd]
		}

		fmt.Printf("MISS line %d @%d\n", lineNo, candStart)
		fmt.Printf("  line: %q\n", visualize(line))
		if nextLine != "" {
			fmt.Printf("  next: %q\n", visualize(nextLine))
		}
	}
}

// helpers

// lastNL возвращает индекс последнего '\n' перед pos, либо -1
func lastNL(s string, pos int) int {
	if pos > len(s) {
		pos = len(s)
	}
	return strings.LastIndexByte(s[:pos], '\n')
}

// nextNL возвращает индекс первого '\n' на/после pos, либо -1
func nextNL(s string, pos int) int {
	if pos < 0 {
		pos = 0
	}
	i := strings.IndexByte(s[pos:], '\n')
	if i < 0 {
		return -1
	}
	return pos + i
}

// visualize — безопасно показывает "невидимые" символы (пробелы, \u00A0, TAB, и т.п.)
func visualize(s string) string {
	var b strings.Builder
	for len(s) > 0 {
		r, sz := utf8.DecodeRuneInString(s)
		switch r {
		case '\n':
			b.WriteString("\\n")
		case '\t':
			b.WriteString("\\t")
		default:
			// Покажем неразрывный пробел и прочие Zs явно
			if r == '\u00A0' {
				b.WriteString("\\u00A0")
			} else if strings.ContainsRune("\u2000\u2001\u2002\u2003\u2004\u2005\u2006\u2007\u2008\u2009\u200A\u202F\u205F\u3000", r) {
				fmt.Fprintf(&b, "\\u%04X", r)
			} else {
				b.WriteRune(r)
			}
		}
		s = s[sz:]
	}
	return b.String()
}
