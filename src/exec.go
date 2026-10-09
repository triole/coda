package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"syscall"

	"github.com/jedib0t/go-pretty/table"
)

func (coda *tCoda) execute(cmds [][]string) (output []byte, exitcode int, err error) {
	var t table.Writer
	if CLI.DryRun {
		t = newTable()
		t.AppendHeader(table.Row{
			"commands that would have been run",
		})
	}

	// Pre-compute tempMap once for all commands (avoid repeated allocations)
	tempMap := coda.makeTempMap(coda.VarMap)

	for _, cmdArr := range cmds {
		cmdArr = coda.iterTemplate(cmdArr, tempMap)
		if CLI.DryRun {
			t.AppendRow(
				[]interface{}{
					fmt.Sprintf("%q", cmdArr),
				},
			)
		} else {
			output, exitcode, err = coda.runCmd(cmdArr)
		}
	}

	if CLI.DryRun {
		fmt.Printf("\n")
		t.Render()
		fmt.Printf("\n")
	}
	return output, exitcode, err
}

func (coda *tCoda) runCmd(cmdArr []string) ([]byte, int, error) {
	var err error
	var exitcode int

	buf := bufferPool.Get().(*bytes.Buffer)
	defer func() { buf.Reset(); bufferPool.Put(buf) }()

	cmd := exec.Command(cmdArr[0], cmdArr[1:]...)
	mw := io.MultiWriter(os.Stdout, buf)

	cmd.Stdout = mw
	cmd.Stderr = mw
	if err = cmd.Run(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				exitcode = status.ExitStatus()
			}
		}
	}
	if err != nil {
		logger.Error("an error occurred: %v", err)
	}

	// Capture output data as a copy before returning buffer to pool
	result := make([]byte, buf.Len())
	copy(result, buf.Bytes())
	return result, exitcode, err
}
