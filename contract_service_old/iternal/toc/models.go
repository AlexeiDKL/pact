package toc

import "regexp"

// TOCItem представляет элемент оглавления
type TOCItem struct {
	Name     string    `json:"name"`
	Caption  string    `json:"caption"`
	StartPos int       `json:"startPosition"`
	EndPos   int       `json:"endPosition"`
	Children []TOCItem `json:"children"`
}

var Header = []string{
	"ДОГОВОР",
	"Приложения",
	"ЧАСТЬ",
	"Приложение",
	"Раздел",
	"Статья",
}

var RegexHeader = map[string]*regexp.Regexp{
	"ДОГОВОР":    regexp.MustCompile(`(?m)^ДОГОВОР\s+\S+\n([^\n]+)`),
	"Приложения": regexp.MustCompile(`(?m)^Приложения\s+\S+\n([^\n]+)`),
	"ЧАСТЬ":      regexp.MustCompile(`(?m)^ЧАСТЬ\s+\S+(?:\s*\n\s*([^\n]+))?`),
	"Приложение": regexp.MustCompile(`(?m)^Приложение\s+N?\s+\d+.`),
	"Раздел":     regexp.MustCompile(`(?m)^Раздел\s+\S+\n([^\n]+)`),
	"Статья":     regexp.MustCompile(`(?m)^Статья\s+\d+\n([^\n]+)`),
}
