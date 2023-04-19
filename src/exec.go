package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"syscall"
)

func (coda tCoda) execute(cmds [][]string) {
	if CLI.Debug {
		fmt.Printf(
			"\n\n%s\n", "Commands (1st template, 2nd what would have been run)",
		)
	}
	for _, cmdArr := range cmds {
		if CLI.Debug {
			fmt.Printf("\n%q\n", cmdArr)
		}
		cmdArr = coda.iterTemplate(cmdArr, coda.VarMap)
		if CLI.Debug {
			fmt.Printf("%q\n", cmdArr)
		} else {
			coda.runCmd(cmdArr)
		}
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
