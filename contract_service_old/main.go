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
	"strings"

	"dkl.ru/pact/contract_service_old/iternal/config"
	httpgarant "dkl.ru/pact/contract_service_old/iternal/http_garant"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

func main() {

	// text, err := readTextFile("D:/work/flutter_project/pact/docs/agree_RU.txt", "utf-8")
	// if err != nil {
	// 	fmt.Println("Ошибка чтения файла:", err)
	// 	return
	// }
	// if text == "" {
	// 	fmt.Println("Пустой текст")
	// 	return
	// }

	// res := toc.TOCItem{
	// 	Name:     "ДОГОВОР",
	// 	Caption:  "ДОГОВОР",
	// 	StartPos: 0,
	// 	EndPos:   toc.ByteIndexAtRune(text, len([]rune(text))),
	// 	Children: []toc.TOCItem{},
	// }

	// // Теперь text содержит содержимое файла
	// items := toc.ParseDocument(text, "", res)

	// // Преобразуем []toc.TOCItem в []toc.TOCItemWithChildren
	// var serializable []toc.TOCItemWithChildren
	// for _, item := range items {
	// 	serializable = append(serializable, toc.ConvertTOCItem(item))
	// }

	// // Оборачиваем в корневой объект с полем Item
	// top := struct {
	// 	Item []toc.TOCItemWithChildren `json:"Item"`
	// }{Item: serializable}

	// // Сохраняем в JSON
	// err = files.SaveToJSON("D:/work/flutter_project/pact/docs/agree_ru_2.json", top)
	// if err != nil {
	// 	fmt.Println("Ошибка сохранения в JSON:", err)
	// 	return
	// }
	download()
}

func download() {
	/*
		70675002
		70675034
		70675048
		70675054
		73184556
		70675082
		70675084
	*/
	err := config.Init()
	if err != nil {
		fmt.Println(err)
		return
	}
	topics := []string{"70675002", "70675034", "70675048", "70675054", "73184556", "70675082", "70675084"}
	for _, v := range topics {
		fmt.Println(v)
		httpgarant.DownloadFromGarantODT(v)
	}
}

func readTextFile(path string, encoding string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("ошибка открытия файла: %w", err)
	}
	defer f.Close()

	var reader io.Reader = f

	switch strings.ToLower(encoding) {
	case "windows-1251", "cp1251":
		reader = transform.NewReader(f, charmap.Windows1251.NewDecoder())
	case "utf-8", "utf8":
		// ничего не делаем, файл уже в utf-8
	default:
		return "", fmt.Errorf("неподдерживаемая кодировка: %s", encoding)
	}

	content, err := io.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("ошибка чтения файла: %w", err)
	}
	return string(content), nil
}
