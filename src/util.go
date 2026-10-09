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
		logger.Fatal("error reading file %q: %v", filename, err)
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
		logger.Fatal("error scanning file %q: %v", filename, err)
	}
	return
}

func makeAbs(filename string) string {
	filename, err := filepath.Abs(filename)
	if err != nil {
		logger.Fatal("can not assemble absolute filename %q", err)
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
		logger.Fatal("unable to determine current user %q", err)
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
