package Utils

import (
	"os"
	"syscall"
)

//////////////////////////////////////////////////////

var UProcesses _Processes_s
type _Processes_s struct {
	/*
		IsPidRunning checks if a process with the given PID is running.

		-----------------------------------------------------------

		> Params:
		  - pid – the PID to check

		> Returns:
		  - true if the process is running, false otherwise
	*/
	IsPidRunning func(pid int) bool
}
//////////////////////////////////////////////////////

func isPidRunningPROCESSES(pid int) bool {
	process, err := os.FindProcess(pid)
	if nil == err {
		err = process.Signal(syscall.Signal(0))

		return nil == err
	} else {
		return false
	}
}
