package files

import "fmt"

//todo доделать и покрыть тестами
func DownloadFromGarantODT(fileId string) error {
	name := "document"
	fileType := "odt"
	filename := fmt.Sprintf("%s.%s", name, fileType)
	url := fmt.Sprintf("https://api.garant.ru/v1/topic/%s/download-odt", fileId)

	fmt.Println(filename, url)
	return nil
}
