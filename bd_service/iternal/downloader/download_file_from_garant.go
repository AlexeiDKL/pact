package downloader

import (
	"fmt"

	"dkl.ru/pact/bd_service/iternal/config"
)

func DownloadFromGarantODT(topicId, filename string) ([]byte, error) {
	url := fmt.Sprintf("https://api.garant.ru/v1/topic/%s/download-odt", topicId)
	token := config.Config.Tokens.Garant

	checksum, err := downloadFile(url, token, filename)
	if err != nil {
		return nil, err
	} else {
		return checksum, nil
	}
}
