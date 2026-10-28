package main

import (
	"os"
)

func (coda *tCoda) SaveFile(data []byte, targetPath string) (err error) {
	tempMap := coda.makeTempMap(coda.VarMap)
	tPath := os.ExpandEnv(coda.execTemplate(targetPath, tempMap))

	log.Info("save file", map[string]any{"path": tPath})
	file, err := os.OpenFile(
		tPath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0600,
	)
	if err != nil {
		log.Fatal("can not open file", map[string]any{"error": err})
	}
	defer file.Close()

	if err == nil {
		_, err = file.Write(data)
		if err != nil {
			log.Fatal("can not write file", map[string]any{"error": err})
		}
	}
	return
}
