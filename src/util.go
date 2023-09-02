package main

import (
	"bufio"
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/jedib0t/go-pretty/table"
	"github.com/jedib0t/go-pretty/text"
)

func getFirstLineOfFile(filename string) (l string) {
	f, err := os.Open(filename)
	if err != nil {
		fmt.Printf("error reading file %q\n", err)
		os.Exit(1)
	}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		l = scanner.Text()
		break
	}
	return
}

func makeAbs(filename string) string {
	filename, err := filepath.Abs(filename)
	if err != nil {
		fmt.Printf("can not assemble absolute filename %q\n", err)
		os.Exit(1)
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
		fmt.Printf("unable to retrieve user's home folder")
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
