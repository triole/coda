package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"sort"
	"strings"
	"sync"
	"text/template"

	"github.com/jedib0t/go-pretty/table"
)

type tVarMap map[string]tVarMapEntry

type tVarMapEntry struct {
	Variable interface{}
	Desc     string
	StrValue string // Cached string value to avoid repeated type assertions
}

// bufferPool reuses bytes.Buffers to reduce allocations
var bufferPool = sync.Pool{
	New: func() interface{} {
		return &bytes.Buffer{}
	},
}

func (vme tVarMapEntry) VarString() string {
	return vme.StrValue
}

func pprint(i interface{}) {
	s, _ := json.MarshalIndent(i, "", "  ")
	fmt.Println(string(s))
}

func makeVarMap(filename string) (varMap tVarMap) {
	varMap = make(tVarMap)

	// Pre-compute values to avoid repeated string operations
	dir := path.Dir(filename)
	base := path.Base(filename)
	ext := strings.TrimPrefix(path.Ext(filename), ".")
	nameNoExt := strings.TrimSuffix(base, "."+ext)
	baseNoExt := strings.TrimSuffix(filename, "."+ext)

	varMap["folder"] = tVarMapEntry{
		Variable: dir, Desc: "folder of file", StrValue: dir,
	}
	varMap["filename"] = tVarMapEntry{Variable: filename, Desc: "full file name", StrValue: filename}
	varMap["shortname"] = tVarMapEntry{Variable: base, Desc: "short name, file name without path", StrValue: base}
	varMap["extension"] = tVarMapEntry{
		Variable: ext, Desc: "file's extension", StrValue: ext,
	}
	varMap["filename_no_ext"] = tVarMapEntry{
		Variable: baseNoExt, Desc: "full file name without preceding extension", StrValue: baseNoExt,
	}
	varMap["shortname_no_ext"] = tVarMapEntry{
		Variable: nameNoExt, Desc: "short name without extension", StrValue: nameNoExt,
	}
	return
}
func (coda *tCoda) makeTempMap(varMap tVarMap) (tempMap map[string]interface{}) {
	// Pre-compute map size for efficiency
	tempMap = make(map[string]interface{}, len(varMap))
	for key, val := range varMap {
		tempMap[key] = val.StrValue // Use cached StrValue directly
	}
	return
}

func (coda *tCoda) iterTemplate(arr []string, varMap tVarMap) (r []string) {
	// Pre-allocate result slice to avoid reallocations
	r = make([]string, len(arr))
	tempMap := coda.makeTempMap(varMap)
	for i, el := range arr {
		r[i] = os.ExpandEnv(coda.execTemplate(el, tempMap))
	}
	return
}

func (coda *tCoda) execTemplate(tplStr string, varMap map[string]interface{}) string {
	// Check cache first (read lock for concurrency)
	coda.cacheMu.RLock()
	tmpl, cached := coda.tmplCache[tplStr]
	coda.cacheMu.RUnlock()

	if !cached {
		// Parse template (outside lock to avoid blocking readers)
		var err error
		tmpl, err = template.New("new.tmpl").Parse(tplStr)
		if err != nil {
			logger.Fatal("template parse error: %w", err)
		}

		// Insert into cache with write lock and double-check
		coda.cacheMu.Lock()
		defer coda.cacheMu.Unlock()

		// Double-check: another goroutine may have inserted it while we parsed
		if tmpl2, ok := coda.tmplCache[tplStr]; ok {
			tmpl = tmpl2
		} else {
			coda.tmplCache[tplStr] = tmpl
		}
	}

	// Use buffer pool to reduce allocations
	buf := bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer bufferPool.Put(buf)

	err := tmpl.Execute(buf, varMap)
	if err != nil {
		logger.Fatal("template execution error: %w", err)
	}
	return buf.String()
}

func orderedIterator(vm tVarMap) (iterator []string) {
	for el := range vm {
		iterator = append(iterator, el)
	}
	sort.Strings(iterator)
	return
}

func printAvailableVars() {
	fmt.Printf("\nAvailable variables\n\n")

	vm := makeVarMap("")
	t := newTable()
	t.AppendHeader(table.Row{
		"variable", "description",
	})
	for _, val := range orderedIterator(vm) {
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
