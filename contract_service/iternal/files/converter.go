package files

import (
	"archive/zip"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
	"unicode"
)

func ConvertOdtToTXT(odtName string) (string, error) {
	txt, err := extractTextFromODT(odtName)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return "", err
	}

	txt = strings.ReplaceAll(txt, "&#160;", " ")
	txt = strings.ReplaceAll(txt, "&quot;", "\"")
	txt = strings.ReplaceAll(txt, "&lt;", "<")
	txt = strings.ReplaceAll(txt, "&gt;", ">")

	return txt, nil
}

func extractTextFromODT(filename string) (string, error) {
	r, err := zip.OpenReader(filename)
	if err != nil {
		return "", err
	}
	defer r.Close()

	var content string
	for _, file := range r.File {
		if file.Name == "content.xml" {
			f, err := file.Open()
			if err != nil {
				return "", err
			}
			defer f.Close()

			data, err := ioutil.ReadAll(f)
			if err != nil {
				return "", err
			}

			content = string(data)
			break
		}
	}

	// Очистка XML от тегов
	re := regexp.MustCompile(`<[^>]+>`)
	cleanText := re.ReplaceAllString(content, "")

	// Очистка пробелов и пустых строк
	cleanText = strings.TrimLeftFunc(cleanText, unicode.IsSpace)

	return cleanText, nil
}
