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
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ztrue/tracerr"

	"VISOR_Server/Utils/Tcef"
)

const (
	// _TEMP_DIR is the full path to the main directory of the temporary files.
	_TEMP_DIR string = _VISOR_DIR + "temp/"
	// _DATA_DIR is the full path to the main directory of the data files.
	_DATA_DIR string = _VISOR_DIR + "data/"
	// _DATA_DIR is the full path to the main directory of the modules.
	_MOD_DIR string = _VISOR_DIR + "Modules/"
)

// _MOD_FOLDER_PREFFIX is the preffix of the modules' folders.
const _MOD_FOLDER_PREFFIX string = "MOD_"

// _MOD_GEN_ERROR_CODE is the exit code of a module when a general error occurs.
const _MOD_GEN_ERROR_CODE int = 3234

const (
	// _MOD_GEN_INFO_JSON is the name of the file containing the module-generated information
	_MOD_GEN_INFO_JSON string = "mod_gen_info.json"
	// _MOD_GEN_INFO_JSON_TMP is the name of the temporary file containing the module-generated information
	_MOD_GEN_INFO_JSON_TMP string = "mod_gen_info.json_tmp"
	// _MOD_USER_INFO_JSON is the name of the file containing the user-given module information (read-only by the
	// module)
	_MOD_USER_INFO_JSON string = "mod_user_info.json"
)

// _MOD_NUMS_NAMES is a map of the numbers of the modules and their names. Use with the NUM_MOD_ constants.
var _MOD_NUMS_NAMES map[int]string = map[int]string{
	2: "S.M.A.R.T. Checker",
	4: "RSS Feed Notifier",
	5: "Email Sender",
}

const (
	NUM_MOD_SMARTChecker    int = 2
	NUM_MOD_RssFeedNotifier int = 4
	NUM_MOD_EmailSender     int = 5
)

// MAX_WAIT_NEXT_TIMESTAMP is the maximum number of seconds to wait for the next timestamp to be registered by a module.
const MAX_WAIT_NEXT_TIMESTAMP int64 = 5

// _RunFileInfo is the struct of the file containing information about the running of a module.
type _RunFileInfo struct {
	// Last_pid is the PID of the last process that ran the module.
	Last_pid int
	// Last_timestamp_ns is the last timestamp in nanoseconds registered by the module.
	Last_timestamp_ns int64
}

// ModGenFileInfo is the struct of the file containing the information about the module.
type ModGenFileInfo[T any] struct {
	// Mod_num is the number of the module.
	Mod_num int
	// Run_info is the information about the running of the module.
	Run_info _RunFileInfo
	// ModSpecificInfo is the information specific to the module, provided by the module - it should be a struct (can be
	// private) and ALL its fields should be exported.
	ModSpecificInfo T
}

// ModProvInfo is the struct that is provided to a module containing information about it.
type ModProvInfo struct {
	// Name is the name of the module.
	Name string
	// Main_Dir is the path to the main directory of the module.
	Main_Dir GPath
	// Data_dir is the path to the directory of the private data files of the module.
	Data_dir GPath
	// Temp_dir is the path to the directory of the private temporary files of the module.
	Temp_dir GPath
}

/*
RealMain is the type of the realMain() function of a module.

realMain is the function that does the actual work of a module (it's equivalent to what main() would normally be).

The generic parameter names are to avoid name conflicts.

-----------------------------------------------------------

– Params:
  - realMain_param_1 – the ModProvInfo struct of the module
  - realMain_param_2 – the ModGenFileInfo struct of the module with the ModGenFileInfo.ModSpecificInfo field of the
    requested type by the module
*/
type RealMain func(realMain_param_1 ModProvInfo, realMain_param_2 any)

/*
ModStartup does the startup routine for a module and executes its realMain() function, catching any fatal errors and
sending an email with them.

Call this as the ONLY thing in the main() function of a module.

-----------------------------------------------------------

– Generic params:
  - T – the type of the ModGenFileInfo.ModSpecificInfo field of the requested type by the module

– Params:
  - mod_num – the number of the module
  - realMain – a pointer to the realMain() function of the module
*/
func ModStartup[T any](mod_num int, realMain RealMain) {
	// Try to run the module, catching any fatal errors and sending an email with them.
	var mod_name string = "ERROR"
	var errors bool = false
	Tcef.Tcef{
		Try: func() {
			// Module startup routine //
			mod_name = GetModNameMODULES(mod_num)
			printStartupSequenceMODULES(mod_name)

			var modGenFileInfo ModGenFileInfo[T] = getModGenFileInfoMODULES[T](mod_num)
			exit, err := processModRunningMODULES(modGenFileInfo)
			if nil != err {
				if err1 := SendModErrorEmailMODULES(mod_num, GetFullErrorMsgGENERAL(err)); nil != err1 {
					fmt.Println("Error sending email with error: " + GetFullErrorMsgGENERAL(err1))
				}
			}
			if exit {
				return
			}

			// Execute realMain()
			realMain(ModProvInfo{
				Name:     mod_name,
				Main_Dir: getModDirMODULES(mod_num),
				Data_dir: getModDataDirMODULES(mod_num),
				Temp_dir: getModTempDirMODULES(mod_num),
			}, modGenFileInfo)
		},
		Catch: func(e Tcef.Exception) {
			errors = true

			var str_email string = ""
			var str_terminal string = ""
			if err, ok := e.(error); ok {
				// tracerr only works with errors
				str_email = GetFullErrorMsgGENERAL(err)
				// Colors for the terminal (not for the email because the colors use ANSI escape codes that are read by
				// the terminal only).
				str_terminal = tracerr.SprintSourceColor(tracerr.Wrap(err), 0)
			} else {
				// If the exception is not an error, get general information about it
				var err_str string = "Invalid type of error information (not a Go \"error\"). " + getVariableInfoGENERAL(e)
				str_email = err_str
				str_terminal = err_str
			}

			// Print the error and send an email with it
			fmt.Println(str_terminal)
			if err := SendModErrorEmailMODULES(mod_num, str_email); nil != err {
				fmt.Println("Error sending email with error: " + GetFullErrorMsgGENERAL(err))
			}
		},
	}.Do()

	// Module shutdown routine //

	if errors {
		printShutdownSequenceMODULES(errors, mod_name, strconv.Itoa(mod_num))

		os.Exit(_MOD_GEN_ERROR_CODE)
	}

	printShutdownSequenceMODULES(errors, mod_name, strconv.Itoa(mod_num))
}

/*
GetModNameMODULES gets the name of a module.

-----------------------------------------------------------

– Params:
  - mod_num – the number of the module

– Returns:
  - the name of the module or an empty string if the module number is invalid
*/
func GetModNameMODULES(mod_num int) string {
	if mod_name, ok := _MOD_NUMS_NAMES[mod_num]; ok {
		return mod_name
	}

	return "INVALID MODULE NUMBER"
}

/*
SendModErrorEmailMODULES directly sends an email to the developer with the error message.

This function does *not* use any modules to do anything. Only utility functions. So it can be used from any
module.

-----------------------------------------------------------

– Params:
  - mod_num – the number of the module from which the error occurred
  - error – the error message

– Returns:
  - nil if the email was sent successfully, otherwise an error
*/
func SendModErrorEmailMODULES(mod_num int, err_str string) error {
	var html_message string = "<pre>" + err_str + "</pre>"

	var html string = *GetModelFileEMAIL(MODEL_FILE_INFO)
	html = strings.ReplaceAll(html, "|3234_HTML_MESSAGE|", html_message)
	html = strings.ReplaceAll(html, "|3234_DATE_TIME|", GetDateTimeStrTIMEDATE())

	return SendEmailEMAIL(prepareEmlEMAIL(EmailInfo{
		Sender:  "VISOR - Info",
		Mail_to: MY_EMAIL_ADDR,
		Subject: "Error in module: " + GetModNameMODULES(mod_num),
		Html:    html,
	}))
}

/*
LoopSleep sleeps for the given number of seconds (with a caveat) and updates the ModGenFileInfo file.

If the number of seconds exceeds MAX_WAIT_NEXT_TIMESTAMP, uses the latter is used instead.

-----------------------------------------------------------

– Params:
  - s – the number of seconds to sleep

– Returns:
  - true if the sleep was successful, false otherwise
*/
func (modGenFileInfo ModGenFileInfo[T]) LoopSleep(s int64) {
	modGenFileInfo.Run_info.Last_timestamp_ns = time.Now().UnixNano()
	modGenFileInfo.Update()

	var seconds = s
	if s > MAX_WAIT_NEXT_TIMESTAMP {
		seconds = MAX_WAIT_NEXT_TIMESTAMP
	}
	time.Sleep(time.Duration(seconds) * time.Second)
}

/*
GetModUserInfo gets the information about the module from the user info file.

-----------------------------------------------------------

– Params:
  - v – a pointer to the variable where the information will be stored, with the struct in which the file is written in

– Returns:
  - true if the file was read successfully, false otherwise
*/
func (modProvInfo ModProvInfo) GetModUserInfo(v any) bool {
	var p_json_file *string = modProvInfo.Data_dir.Add(_MOD_USER_INFO_JSON).ReadFile()
	if p_json_file == nil {
		return false
	}

	return FromJsonGENERAL([]byte(*p_json_file), v)
}

/*
Update updates the information about the module in its generated information file.

-----------------------------------------------------------

– Params:
  - mod_num – the number of the module

– Returns:
  - nil if the update was successful, false otherwise
*/
func (modGenFileInfo ModGenFileInfo[T]) Update() error {
	var json_str string = *ToJsonGENERAL(&modGenFileInfo)

	var file_path_curr GPath = getModDataDirMODULES(modGenFileInfo.Mod_num).Add(_MOD_GEN_INFO_JSON)
	var file_path_new GPath = getModDataDirMODULES(modGenFileInfo.Mod_num).Add(_MOD_GEN_INFO_JSON_TMP)

	var err error = file_path_new.WriteTextFile(json_str)
	if nil != err {
		return err
	}

	return os.Rename(file_path_new.GPathToStringConversion(), file_path_curr.GPathToStringConversion())
}

func printStartupSequenceMODULES(mod_name string) {
	fmt.Println("//------------------------------------------\\\\")
	fmt.Println("--- " + mod_name + " ---")
	fmt.Println("V.I.S.O.R. Systems")
	fmt.Println("------------------")
	fmt.Println()
}

func printShutdownSequenceMODULES(errors bool, mod_name string, mod_num string) {
	fmt.Println()
	fmt.Println("---------")
	if errors {
		fmt.Println("Exiting with ERRORS the module \"" + mod_name + "\" (number " + mod_num + ")...")
	} else {
		fmt.Println("Exiting normally the module \"" + mod_name + "\" (number " + mod_num + ")...")
	}
	fmt.Println("\\\\------------------------------------------//")
}

/*
getModDirMODULES gets the full path to the directory of a module.

-----------------------------------------------------------

– Params:
  - mod_num – the number of the module

– Returns:
  - the full path to the directory of the module
*/
func getModDirMODULES(mod_num int) GPath {
	return PathFILESDIRS(_MOD_DIR, _MOD_FOLDER_PREFFIX + strconv.Itoa(mod_num) + "/")
}

/*
getModDataDirMODULES gets the full path to the private data directory of a module.

-----------------------------------------------------------

– Params:
  - mod_num – the number of the module

– Returns:
  - the full path to the private data directory of the module
*/
func getModDataDirMODULES(mod_num int) GPath {
	return PathFILESDIRS(_DATA_DIR, _MOD_FOLDER_PREFFIX + strconv.Itoa(mod_num) + "/")
}

/*
getModTempDirMODULES gets the full path to the private temporary directory of a module.

-----------------------------------------------------------

– Params:
  - mod_num – the number of the module

– Returns:
  - the full path to the private temporary directory of the module
*/
func getModTempDirMODULES(mod_num int) GPath {
	return PathFILESDIRS(_TEMP_DIR, _MOD_FOLDER_PREFFIX + strconv.Itoa(mod_num) + "/")
}

/*
processModRunningMODULES checks if the module is already running and exits if it is, and if it's not, writes the
necessary information to the module info file.

-----------------------------------------------------------

– Params:
  - modGenFileInfo – the information of the module

– Returns:
  - true if the module is already running, false otherwise
  - nil if the module information was updated, false otherwise
*/
func processModRunningMODULES[T any](modGenFileInfo ModGenFileInfo[T]) (bool, error) {
	// Check PID and timestamp
	if modGenFileInfo.Run_info.Last_pid != -1 && IsPidRunningPROCESSES(modGenFileInfo.Run_info.Last_pid) &&
		(time.Now().UnixNano() - modGenFileInfo.Run_info.Last_timestamp_ns) < (MAX_WAIT_NEXT_TIMESTAMP * 1e9) {

		// todo This is temporary, to see if the modules are being started many times in a row almost instantaneously
		panic("Module already running")

		return true, nil
	}

	modGenFileInfo.Run_info.Last_pid = os.Getpid()
	modGenFileInfo.Run_info.Last_timestamp_ns = time.Now().UnixNano()

	return false, modGenFileInfo.Update()
}

/*
getModGenFileInfoMODULES gets the information of a module from the module info file or creates a new one if it doesn't
exist.

-----------------------------------------------------------

– Params:
  - mod_num – the number of the module

– Returns:
  - the information of the module
*/
func getModGenFileInfoMODULES[T any](mod_num int) ModGenFileInfo[T] {
	var info ModGenFileInfo[T]

	// Check first if the temporary file exists
	var p_info *string = getModDataDirMODULES(mod_num).Add(_MOD_GEN_INFO_JSON_TMP).ReadFile()
	if nil == p_info {
		// If not, check if the main file exists
		p_info = getModDataDirMODULES(mod_num).Add(_MOD_GEN_INFO_JSON).ReadFile()
		if nil == p_info {
			// If not, write a new file

			goto write_file
		}
	}

	if FromJsonGENERAL([]byte(*p_info), &info) {
		return info
	}

write_file:

	info.Mod_num = mod_num
	info.Run_info.Last_pid = -1
	info.Run_info.Last_timestamp_ns = -1
	info.Update()

	return info
}
