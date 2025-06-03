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
	"Договор",
	"Приложения",
	"Часть",
	"Приложение",
	"Раздел",
	"Статья",
}

var RegexHeader = map[string]*regexp.Regexp{
	"Договор":    regexp.MustCompile(`(?m)^Договор\s+\S+\n([^\n]+)`),
	"Приложения": regexp.MustCompile(`(?m)^Приложения\s+\S+\n([^\n]+)`),
	"Часть":      regexp.MustCompile(`(?m)^Часть\s+\S+\n([^\n]+)`),
	"Приложение": regexp.MustCompile(`(?m)^Приложение\s+\S+\n([^\n]+)`),
	"Раздел":     regexp.MustCompile(`(?m)^Раздел\s+\S+\n([^\n]+)`),
	"Статья":     regexp.MustCompile(`(?m)^Статья\s+\d+\n([^\n]+)`),
}
