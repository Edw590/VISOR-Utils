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
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"Utils/Tcef"
)

const (
	// _DATA_REL_DIR is the relative path to the data directory from PersonalConsts._VISOR_DIR.
	_DATA_REL_DIR string = "data/"
	// _TEMP_FOLDER is the relative path to the temporary folder from PersonalConsts._VISOR_DIR.
	_TEMP_FOLDER string = _DATA_REL_DIR + "Temp/"
	// _USER_DATA_REL_DIR is the relative path to the user data directory from PersonalConsts._VISOR_DIR.
	_USER_DATA_REL_DIR string = _DATA_REL_DIR + "UserData/"
	// _PROGRAM_DATA_REL_DIR is the relative path to the program data directory from PersonalConsts._VISOR_DIR.
	_PROGRAM_DATA_REL_DIR string = _DATA_REL_DIR + "ProgramData/"
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
	// ProgramData_dir is the path to the directory of the program data files.
	ProgramData_dir GPath
	// UserData_dir is the path to the directory of the private user data files.
	UserData_dir GPath
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
	var errs bool = false
	Tcef.Tcef{
		Try: func() {
			// Initialize the personal "constants"
			err := PersonalConsts_GL.init()
			if err != nil {
				fmt.Println("CRITICAL ERROR: " + GetFullErrorMsgGENERAL(err))
				errs = true

				return
			}

			// Module startup routine //
			mod_name = GetModNameMODULES(mod_num)
			printStartupSequenceMODULES(mod_name)

			exit, err, modGenFileInfo := processModRunningMODULES[T](mod_num)
			if nil != err {
				var str_error string = GetFullErrorMsgGENERAL(err)
				if err1 := SendModErrorEmailMODULES(mod_num,str_error); nil != err1 {
					fmt.Println("Error sending email with errors\n" +
						GetFullErrorMsgGENERAL(err1) + "\n-----\n" + str_error)
				}
			}
			if exit {
				return
			}

			// Execute realMain()
			realMain(ModProvInfo{
				Name:            mod_name,
				ProgramData_dir: getProgramDataDirMODULES(mod_num),
				UserData_dir:    getUserDataDirMODULES(mod_num),
				Temp_dir:        getModTempDirMODULES(mod_num),
			}, modGenFileInfo)
		},
		Catch: func(e Tcef.Exception) {
			errs = true

			var str_error string = GetFullErrorMsgGENERAL(e)

			// Print the error and send an email with it
			fmt.Println(str_error)
			if err := SendModErrorEmailMODULES(mod_num, str_error); nil != err {
				fmt.Println("Error sending email with error:\n" + GetFullErrorMsgGENERAL(err) + "\n-----\n" + str_error)
			}
		},
	}.Do()

	// Module shutdown routine //

	if errs {
		printShutdownSequenceMODULES(errs, mod_name, strconv.Itoa(mod_num))

		os.Exit(_MOD_GEN_ERROR_CODE)
	}

	printShutdownSequenceMODULES(errs, mod_name, strconv.Itoa(mod_num))
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

	var p_html *string = GetModelFileEMAIL(MODEL_FILE_INFO)
	if nil == p_html {
		return errors.New("error getting email model file \"" + MODEL_FILE_INFO + "\"")
	}
	var html string = *p_html
	html = strings.ReplaceAll(html, "|3234_HTML_MESSAGE|", html_message)
	html = strings.ReplaceAll(html, "|3234_DATE_TIME|", GetDateTimeStrTIMEDATE())

	message_eml, mail_to, success := prepareEmlEMAIL(EmailInfo{
		Sender:  "VISOR - Info",
		Mail_to: PersonalConsts_GL.USER_EMAIL_ADDR,
		Subject: "Error in module: " + GetModNameMODULES(mod_num),
		Html:    html,
	})
	if !success {
		return errors.New("error preparing email")
	}

	return SendEmailEMAIL(message_eml, mail_to, true)
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
func (modGenFileInfo ModGenFileInfo[T]) LoopSleep(s int64) error {
	modGenFileInfo.Run_info.Last_timestamp_ns = time.Now().UnixNano()
	var err error = modGenFileInfo.Update()

	var seconds = s
	if s > MAX_WAIT_NEXT_TIMESTAMP {
		seconds = MAX_WAIT_NEXT_TIMESTAMP
	}
	time.Sleep(time.Duration(seconds) * time.Second)

	return err
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
	var p_json_file *string = modProvInfo.UserData_dir.Add(_MOD_USER_INFO_JSON).ReadFile()
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

	var file_path_curr GPath = getUserDataDirMODULES(modGenFileInfo.Mod_num).Add(_MOD_GEN_INFO_JSON)
	var file_path_new GPath = getUserDataDirMODULES(modGenFileInfo.Mod_num).Add(_MOD_GEN_INFO_JSON_TMP)

	var err error = file_path_new.WriteTextFile(json_str)
	if nil != err {
		return err
	}

	return os.Rename(file_path_new.GPathToStringConversion(), file_path_curr.GPathToStringConversion())
}

/*
printStartupSequenceMODULES prints the startup sequence of a module.

-----------------------------------------------------------

– Params:
  - mod_name – the name of the module
*/
func printStartupSequenceMODULES(mod_name string) {
	fmt.Println("//------------------------------------------\\\\")
	fmt.Println("--- " + mod_name + " ---")
	fmt.Println("V.I.S.O.R. Systems")
	fmt.Println("------------------")
	fmt.Println()
}

/*
printShutdownSequenceMODULES prints the shutdown sequence of a module.

-----------------------------------------------------------

– Params:
  - errors – true if the module is exiting with errors, false otherwise
  - mod_name – the name of the module
  - mod_num – the number of the module
*/
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
getProgramDataDirMODULES gets the full path to the program data directory of a module.

-----------------------------------------------------------

– Params:
  - mod_num – the number of the module

– Returns:
  - the full path to the program data directory of the module
*/
func getProgramDataDirMODULES(mod_num int) GPath {
	return PathFILESDIRS(PersonalConsts_GL._VISOR_DIR, _PROGRAM_DATA_REL_DIR, _MOD_FOLDER_PREFFIX + strconv.Itoa(mod_num) + "/")
}

/*
getUserDataDirMODULES gets the full path to the private user data directory of a module.

-----------------------------------------------------------

– Params:
  - mod_num – the number of the module

– Returns:
  - the full path to the private data directory of the module
*/
func getUserDataDirMODULES(mod_num int) GPath {
	return PathFILESDIRS(PersonalConsts_GL._VISOR_DIR, _USER_DATA_REL_DIR, _MOD_FOLDER_PREFFIX + strconv.Itoa(mod_num) + "/")
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
	return PathFILESDIRS(PersonalConsts_GL._VISOR_DIR, _TEMP_FOLDER, _MOD_FOLDER_PREFFIX + strconv.Itoa(mod_num) + "/")
}

/*
processModRunningMODULES checks if the module is already running and exits if it is, and writes the necessary
information to the module info files.

-----------------------------------------------------------

– Params:
  - mod_num – the number of the module

– Returns:
  - true if the module should stop running for any reason, false if it should carry on (already running, no directory
	access, etc.)
  - nil if the module information was updated, false otherwise
  - the information of the module
*/
func processModRunningMODULES[T any](mod_num int) (bool, error, ModGenFileInfo[T]) {
	var curr_pid int = os.Getpid()
	var curr_ts int64 = time.Now().UnixNano()

	if err := getUserDataDirMODULES(mod_num).Add(
				"PID=" + strconv.Itoa(curr_pid) +
				"_TS=" + strconv.FormatInt(curr_ts, 10),
			).Create(true); nil != err {
		return true, err, ModGenFileInfo[T]{}
	}

	files, err := os.ReadDir(getUserDataDirMODULES(mod_num).GPathToStringConversion())
	if nil != err {
		return true, err, ModGenFileInfo[T]{}
	}

	for _, file := range files {
		if strings.HasPrefix(file.Name(), "PID=") {
			var file_path GPath = getUserDataDirMODULES(mod_num).Add(file.Name())

			var info_list []string = strings.Split(file.Name(), "_")
			var pid_str string = strings.TrimPrefix(info_list[0], "PID=")
			var ts_str string = strings.TrimPrefix(info_list[1], "TS=")

			var pid int
			if pid, err = strconv.Atoi(pid_str); nil != err {
				_ = file_path.Remove()

				continue
			}
			var ts int64
			if ts, err = strconv.ParseInt(ts_str, 10, 64); nil != err {
				_ = file_path.Remove()

				continue
			}

			if pid != curr_pid {
				if IsPidRunningPROCESSES(pid) && (curr_ts - ts) < (MAX_WAIT_NEXT_TIMESTAMP*1e9) {

					// todo This is temporary, to see when the modules are being started many times in a row almost instantaneously
					PanicGENERAL("Module already running")

					_ = file_path.Remove()

					return true, nil, ModGenFileInfo[T]{}
				}
			}
		}
	}

	for _, file := range files {
		if strings.HasPrefix(file.Name(), "PID=") {
			var file_path GPath = getUserDataDirMODULES(mod_num).Add(file.Name())

			var info_list []string = strings.Split(file.Name(), "_")
			var pid_str string = strings.TrimPrefix(info_list[0], "PID=")

			var pid int
			if pid, err = strconv.Atoi(pid_str); nil != err || pid != curr_pid {
				_ = file_path.Remove()
			}
		}
	}

	var modGenFileInfo ModGenFileInfo[T]

	// Check first if the temporary file exists
	var p_info *string = getUserDataDirMODULES(mod_num).Add(_MOD_GEN_INFO_JSON_TMP).ReadFile()
	if nil == p_info {
		// If not, check if the main file exists
		p_info = getUserDataDirMODULES(mod_num).Add(_MOD_GEN_INFO_JSON).ReadFile()
		if nil == p_info {
			// If not, write a new file

			goto new_file
		}
	}

	FromJsonGENERAL([]byte(*p_info), &modGenFileInfo)

new_file:

	modGenFileInfo.Mod_num = mod_num
	modGenFileInfo.Run_info.Last_pid = curr_pid
	modGenFileInfo.Run_info.Last_timestamp_ns = curr_ts

	return false, modGenFileInfo.Update(), modGenFileInfo
}
