package main

import (
	"fmt"

	files "dkl.ru/pact/contract_service/iternal/files"
)

func main() {
	fmt.Println("contract_service")
	files.DownloadFromGarantODT("123")
}
