package main

/*
	проверяет обновления договора, 	httpgarant.CheckFileUpdate
	загружает новый текст, 			httpgarant.DownloadFromGarantODT
	конвертирует в txt, 			files.ConvertOdtToTXT
	формирует оглавление			toc.ParseDocument
*/

import (
	"fmt"

	httpgarant "dkl.ru/pact/contract_service_old/iternal/http_garant"
	"dkl.ru/pact/contract_service_old/iternal/initialization"
	"dkl.ru/pact/contract_service_old/iternal/logger"
)

func main() {
	err := initialization.Init()
	if err != nil {
		panic(err.Error())
	}

	varr, err := httpgarant.CheckFileUpdate(70670880, "2025-04-29")
	if err != nil {
		logger.Logger.Info(fmt.Sprintf("%s!", err.Error()))
	}
	logger.Logger.Debug(fmt.Sprintf("%t", varr))
}
