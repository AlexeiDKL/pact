package main

/*
	проверяет обновления договора, 	httpgarant.CheckFileUpdate
	загружает новый текст, 			httpgarant.DownloadFromGarantODT
	конвертирует в txt, 			files.ConvertOdtToTXT
	формирует оглавление			toc.ParseDocument
*/

import (
	"fmt"
	"io/ioutil"

	"dkl.ru/pact/contract_service_old/iternal/toc"
)

func main() {
	textPath := "agree_ru.txt"

	// Чтение файла в переменную
	content, err := ioutil.ReadFile(textPath)
	if err != nil {
		fmt.Println("Ошибка чтения файла:", err)
		return
	}
	text := string(content) // если нужен текст

	res := toc.TOCItem{
		Name:     "ДОГОВОР",
		Caption:  "ДОГОВОР",
		StartPos: 0,
		EndPos:   len(text),
		Children: []toc.TOCItem{},
	}

	// Теперь text содержит содержимое файла
	items := toc.ParseDocument(text, "", res)
	fmt.Println(items)

}
