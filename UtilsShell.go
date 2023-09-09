/*******************************************************************************
 * Copyright 2023-2023 Edw590
 *
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 ******************************************************************************/

package Utils

import (
	"bytes"
	"os/exec"
	"runtime"
	"strings"
)

// StdOutErrCmd is a struct containing the stdout and stderr of a command.
type StdOutErrCmd struct {
	// Stdout_str is the stdout of the command as a string.
	Stdout_str string
	// Stdout is the stdout of the command as a buffer.
	Stdout *bytes.Buffer
	// Stderr_str is the stderr of the command as a string.
	Stderr_str string
	// Stderr is the stderr of the command as a buffer.
	Stderr *bytes.Buffer
}

/*
ExecCmdSHELL executes a command in a shell and returns the stdout and stderr.

On Windows, the command is executed in powershell.exe; on Linux, it's executed in bash.

ATTENTION: to call any program, add "{{EXE}}" right after the program name in the command string. This will be replaced
by ".exe" on Windows and by "" on Linux. This avoids PowerShell aliases ("curl" is an alias for "Invoke-WebRequest" in
PowerShell but the actual program in Linux, for example - but curl.exe calls the actual program).

-----------------------------------------------------------

– Params:
  - command – the command to execute

– Returns:
  - the StdOutErrCmd struct containing the stdout and stderr of the command. Note that their string versions have all
    line endings replaced with "\n".
  - the error returned by the command execution, if any
*/
func ExecCmdSHELL(command string) (StdOutErrCmd, error) {
	var commands []string = nil
	if "windows" == runtime.GOOS {
		commands = []string{"powershell.exe", strings.Replace(command, "{{EXE}}", ".exe", -1)}
	} else {
		commands = []string{"bash", "-c", strings.Replace(command, "{{EXE}}", "", -1)}
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command(commands[0], commands[1:]...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()

	var stdout_str = strings.ReplaceAll(stdout.String(), "\r\n", "\n")
	stdout_str = strings.ReplaceAll(stdout_str, "\r", "\n")
	var stderr_str = strings.ReplaceAll(stderr.String(), "\r\n", "\n")
	stderr_str = strings.ReplaceAll(stderr_str, "\r", "\n")

	return StdOutErrCmd{
		Stdout_str: stdout_str,
		Stdout: &stdout,
		Stderr_str: stderr_str,
		Stderr: &stderr,
	}, err
}
