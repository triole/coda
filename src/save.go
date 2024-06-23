package main

import (
	"fmt"
	"log"
	"os"
)

func (coda tCoda) SaveFile(data []byte, targetPath string) (err error) {
	tempMap := coda.makeTempMap(coda.VarMap)
	tPath := os.ExpandEnv(coda.execTemplate(targetPath, tempMap))

	fmt.Printf("[coda] save file %q\n", tPath)
	file, err := os.OpenFile(
		tPath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666,
	)
	if err != nil {
		log.Fatal("[coda] can not open file: ", err)
	}
	defer file.Close()

	if err == nil {
		_, err = file.Write(data)
		if err != nil {
			log.Fatal("[coda] can not write file: ", err)
		}
	}
	return
}
