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
	"ДОГОВОР":    regexp.MustCompile(`(?m)^ДОГОВОР`),
	"Приложения": regexp.MustCompile(`(?m)^Приложения`),
	"ЧАСТЬ":      regexp.MustCompile(`(?m)^ЧАСТЬ\s+\S+`),
	// "Приложение": regexp.MustCompile(`(?im)^[ \t]*приложение[ \t]*(?:№[ \t]*)?([A-ZА-ЯЁIVXLCDM0-9]+)?[ \t]*$`),
	"Приложение": regexp.MustCompile(`(?im)^[\p{Zs}\t]*(приложение[\p{Zs}\t]*(?:(?:№|N)[\p{Zs}\t]*\d+|\d+)?)[\p{Zs}\t]*[^\r\n]*`),
	"Раздел":     regexp.MustCompile(`(?m)^Раздел\s+[IVXLCDM]+`),
	"Статья":     regexp.MustCompile(`(?m)^Статья\s+\d+`),
}
