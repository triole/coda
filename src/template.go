package main

import (
	"bytes"
	"fmt"
	"os"
	"sort"
	"strings"
	"text/template"

	"github.com/jedib0t/go-pretty/table"
	"github.com/jedib0t/go-pretty/text"
)

type tVarMap map[string]tVarMapEntry

type tVarMapEntry struct {
	Variable interface{}
	Desc     string
}

func (vme tVarMapEntry) VarString() string {
	switch val := vme.Variable.(type) {
	case string:
		return val
	default:
		return ""
	}
}

func makeVarMap(filename string) (varMap tVarMap) {
	varMap = make(tVarMap)
	varMap["filename"] = tVarMapEntry{filename, "full file name"}
	varMap["shortname"] = tVarMapEntry{
		find(`[^/]+$`, filename), "short name, file name without path",
	}
	ext := ""
	varMap["extension"] = tVarMapEntry{"", "file's extension"}
	if strings.Contains(filename, ".") == true {
		arr := strings.Split(filename, ".")
		if len(arr)-1 > 0 {
			ext = arr[len(arr)-1]
		}
		varMap["extension"] = tVarMapEntry{ext, "extension of file"}
	}
	varMap["filename_no_ext"] = tVarMapEntry{
		strings.Replace(
			filename, "."+varMap["extension"].VarString(), "", -1,
		), "full file name without preceeding extension",
	}
	varMap["shortname_no_ext"] = tVarMapEntry{
		strings.Replace(
			varMap["shortname"].VarString(), "."+varMap["ext"].VarString(), "", -1,
		), "short name without extension",
	}
	return
}

func (coda tCoda) iterTemplate(arr []string, varMap tVarMap) (r []string) {
	tempMap := make(map[string]interface{})
	for key, val := range varMap {
		tempMap[key] = val.VarString()
	}
	for _, el := range arr {
		r = append(r, coda.execTemplate(el, tempMap))
	}
	return
}

func (coda tCoda) execTemplate(tplStr string, varMap map[string]interface{}) string {
	tmpl := template.Must(
		template.New("new.tmpl").Parse(tplStr),
	)

	buf := &bytes.Buffer{}
	err := tmpl.Execute(buf, varMap)
	if err != nil {
		panic(err)
	}
	return buf.String()
}

func printAvailableVars() {
	vm := makeVarMap("")
	var iterator []string
	for el := range vm {
		iterator = append(iterator, el)
	}
	sort.Strings(iterator)

	fmt.Printf("\nAvailable variables\n\n")
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
			Header: text.FormatUpper,
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
	t.AppendHeader(table.Row{
		"variable", "description",
	})
	for _, val := range iterator {
		t.AppendRow(
			[]interface{}{
				"{{." + val + "}}",
				vm[val].Desc,
			},
		)
	}
	t.Render()
	fmt.Printf("\n")
}
