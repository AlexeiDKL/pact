package garant

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"dkl.ru/pact/garant_service/iternal/config"
	myerrors "dkl.ru/pact/garant_service/iternal/my_errors"
)

func DownloadODT(topic, filePath string) error {
	url := fmt.Sprintf("https://api.garant.ru/v1/topic/%s/download-odt", topic)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	fmt.Println(url, config.Config.Tokens.Garant, "!")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.Config.Tokens.Garant)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return myerrors.NotDownload(filePath, fmt.Errorf("status code: %d", resp.StatusCode))
	}

	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
