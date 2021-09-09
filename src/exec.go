package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func (coda tCoda) execute(cmds [][]string) {
	for _, cmdArr := range cmds {
		cmdArr = coda.expandVars(cmdArr)
		if CLI.Debug == true {
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

func (coda tCoda) expandVars(cmdArr []string) (r []string) {
	r = cmdArr
	for idx, el := range cmdArr {
		for key, val := range coda.VarMap {
			new := strings.Replace(
				el, "{"+strings.ToUpper(key)+"}", val, -1,
			)
			if new != el {
				r[idx] = new
			}
		}
	}
	return cmdArr
}
