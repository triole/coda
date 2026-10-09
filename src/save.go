package main

import (
	"os"
)

func (coda *tCoda) SaveFile(data []byte, targetPath string) (err error) {
	tempMap := coda.makeTempMap(coda.VarMap)
	tPath := os.ExpandEnv(coda.execTemplate(targetPath, tempMap))

	logger.Info("save file %q", tPath)
	file, err := os.OpenFile(
		tPath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0600,
	)
	if err != nil {
		logger.Fatal("can not open file: ", err)
	}
	defer file.Close()

	if err == nil {
		_, err = file.Write(data)
		if err != nil {
			logger.Fatal("can not write file: ", err)
		}
	}
	return
}
