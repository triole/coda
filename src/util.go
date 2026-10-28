package main

import (
	"bufio"
	"os"
	"os/user"
	"path/filepath"

	"github.com/jedib0t/go-pretty/table"
	"github.com/jedib0t/go-pretty/text"
)

func getFirstLineOfFile(filename string) (l string) {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal("error reading file", map[string]any{"file": filename, "error": err})
	}
	scanner := bufio.NewScanner(f)
	// Limit line size to prevent OOM on malicious/extremely long lines
	buf := make([]byte, 64*1024)
	scanner.Buffer(buf, 1024*1024)
	for scanner.Scan() {
		l = scanner.Text()
		break
	}
	if err := scanner.Err(); err != nil {
		log.Fatal("error scanning file", map[string]any{"file": filename, "error": err})
	}
	return
}

func makeAbs(filename string) string {
	filename, err := filepath.Abs(filename)
	if err != nil {
		log.Fatal("can not assemble absolute filename", map[string]any{"error": err})
	}
	return filename
}

func isFile(filePath string) bool {
	stat, err := os.Stat(makeAbs(filePath))
	if !os.IsNotExist(err) && !stat.IsDir() {
		return true
	}
	return false
}

func getHome() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal("unable to determine current user", map[string]any{"error": err})
	}
	return usr.HomeDir
}

func newTable() table.Writer {
	t := table.NewWriter()

	t.SetStyle(table.Style{
		Name: "myNewStyle",
		Box: table.BoxStyle{
			MiddleHorizontal: "-",
			MiddleSeparator:  "+",
			MiddleVertical:   "|",
			PaddingLeft:      " ",
			PaddingRight:     " ",
		},
		Format: table.FormatOptions{
			Header: text.FormatLower,
		},
		Options: table.Options{
			DrawBorder:      false,
			SeparateColumns: true,
			SeparateFooter:  true,
			SeparateHeader:  true,
			SeparateRows:    false,
		},
	})

	t.SetOutputMirror(os.Stdout)
	return t
}
