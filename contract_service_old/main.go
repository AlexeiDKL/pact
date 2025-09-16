package main

/*
	проверяет обновления договора, 	httpgarant.CheckFileUpdate
	загружает новый текст, 			httpgarant.DownloadFromGarantODT
	конвертирует в txt, 			files.ConvertOdtToTXT
	формирует оглавление			toc.ParseDocument
*/

import (
	"fmt"
	"io"
	"os"

	"dkl.ru/pact/contract_service_old/iternal/files"
	"dkl.ru/pact/contract_service_old/iternal/toc"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

func main() {

	text := getText()
	if text == "" {
		fmt.Println("Пустой текст")
		return
	}

	res := toc.TOCItem{
		Name:     "ДОГОВОР",
		Caption:  "ДОГОВОР",
		StartPos: 0,
		EndPos:   len(text) - 1,
		Children: []toc.TOCItem{},
	}

	// Теперь text содержит содержимое файла
	items := toc.ParseDocument(text, "", res)

	// Преобразуем []toc.TOCItem в []toc.TOCItemWithChildren
	var serializable []toc.TOCItemWithChildren
	for _, item := range items {
		serializable = append(serializable, toc.ConvertTOCItem(item))
	}

	// Оборачиваем в корневой объект с полем Item
	top := struct {
		Item []toc.TOCItemWithChildren `json:"Item"`
	}{Item: serializable}

	// Сохраняем в JSON
	err := files.SaveToJSON("files/agree_ru_2.json", top)
	if err != nil {
		fmt.Println("Ошибка сохранения в JSON:", err)
		return
	}
}

func getText() string {
	textPath := "files/agree_ru.txt"

	// Чтение файла в Windows-1251 и декодирование в UTF-8
	f, err := os.Open(textPath)
	if err != nil {
		fmt.Println("Ошибка открытия файла:", err)
		return ""
	}
	defer f.Close()

	reader := transform.NewReader(f, charmap.Windows1251.NewDecoder())
	content, err := io.ReadAll(reader)
	if err != nil {
		fmt.Println("Ошибка чтения файла:", err)
		return ""
	}
	text := string(content)
	return text
}
