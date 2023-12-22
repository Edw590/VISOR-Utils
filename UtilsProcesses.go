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
	"os"
	"os/exec"
	"runtime"
	"syscall"
)

/*
IsPidRunningPROCESSES checks if a process with the given PID is running.

-----------------------------------------------------------

– Params:
  - pid – the PID to check

– Returns:
  - true if the process is running, false otherwise
*/
func IsPidRunningPROCESSES(pid int) bool {
	if pid < 0 {
		return false
	}

	process, err := os.FindProcess(pid)
	if nil == err {
		if runtime.GOOS == "windows" {
			return true
		}

		err = process.Signal(syscall.Signal(0))

		return nil == err
	} else {
		return false
	}
}

/*
StartProcessPROCESSES starts a new separate process with the given path.

-----------------------------------------------------------

– Params:
  - path – the path of the program to start

– Returns:
  - true if the process was started, false otherwise
 */
func StartProcessPROCESSES(path GPath) bool {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("powershell.exe", "/C", "start", path.GPathToStringConversion())
		err := cmd.Start()
		if err != nil {
			return false
		}
		err = cmd.Process.Release()
		if err != nil {
			return false
		}
	} else {
		cmd := exec.Command("sh", "-c", path.GPathToStringConversion(), "&")
		err := cmd.Start()
		if err != nil {
			return false
		}
		err = cmd.Process.Release()
		if err != nil {
			return false
		}
	}

	return true
}
