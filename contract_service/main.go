package main

import (
	files "dkl.ru/pact/contract_service/iternal/files"
	"dkl.ru/pact/contract_service/iternal/initialization"
	"dkl.ru/pact/contract_service/iternal/logger"
)

func main() {
	err := initialization.Init()
	if err != nil {
		panic(err.Error())
	}
	err = files.DownloadFromGarantODT("123")
	logger.Logger.Info(err.Error())
}
