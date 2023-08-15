package Utils

import (
	"bytes"
	"os/exec"
	"runtime"
)

//////////////////////////////////////////////////////

var UShell _Shell_s
type _Shell_s struct {
	/*
		ExecCmd executes a command in the shell and returns the stdout and stderr.

		-----------------------------------------------------------

		> Params:
		  - command – the command to execute

		> Returns:
		  - the stdout
		  - the stderr
		  - the error of running the command
	*/
	ExecCmd func(command string) (string, string, error)
}
//////////////////////////////////////////////////////

func execCmdSHELL(command string) (string, string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	var shell_to_use string = ""
	if "windows" == runtime.GOOS {
		shell_to_use = "cmd"
	} else {
		shell_to_use = "bash"
	}

	cmd := exec.Command(shell_to_use, "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()

	return stdout.String(), stderr.String(), err
}
