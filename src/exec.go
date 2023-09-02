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

func (coda tCoda) execute(cmds [][]string) {
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
			coda.runCmd(cmdArr)
		}
	}

	if CLI.DryRun {
		fmt.Printf("\n")
		t.Render()
		fmt.Printf("\n")
	}
}

func (coda tCoda) runCmd(cmdArr []string) ([]byte, int, error) {
	var err error
	var exitcode int
	var stdBuffer bytes.Buffer

	cmd := exec.Command(cmdArr[0], cmdArr[1:]...)
	// mw := io.MultiWriter(&stdBuffer)
	mw := io.MultiWriter(os.Stdout, &stdBuffer)

	cmd.Stdout = mw
	cmd.Stderr = mw
	if err = cmd.Run(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			// the program has exited with an exit code != 0
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				exitcode = status.ExitStatus()
			}
		}
	}
	if err != nil {
		fmt.Printf("An error occured: %s\n", err)
	}
	fmt.Printf("")
	return stdBuffer.Bytes(), exitcode, err
}
