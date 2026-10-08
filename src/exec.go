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
	for _, cmdArr := range cmds {
		cmdArr = coda.iterTemplate(cmdArr, coda.VarMap)
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

	// Use buffer pool to reduce allocations
	buf := bufferPool.Get().(*bytes.Buffer)
	defer bufferPool.Put(buf)

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
		logger.Error("an error occured: %q\n", err)
	}
	return buf.Bytes(), exitcode, err
}
