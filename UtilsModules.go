package Utils

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ztrue/tracerr"

	PersonalConsts "VISOR_S_L/PERSONAL_FILES_EOG"
	"VISOR_S_L/Utils/Tcf"
)

//////////////////////////////////////////////////////

var UModules _Modules_s
type _Modules_s struct {
	/*
		GetModName gets the name of a module.

		-----------------------------------------------------------

		> Params:
		  - mod_num – the number of the module

		> Returns:
		  - the name of the module or an empty string if the module number is invalid
	*/
	GetModName func(mod_num int) string
	/*
		SendModErrorEmail directly sends an email to the developer with the error message.

		This function does *not* use any modules to do anything. Only utility functions. So it can be used from any
		module.

		-----------------------------------------------------------

		> Params:
		  - mod_num – the number of the module from which the error occurred
		  - error – the error message

		> Returns:
		  - true if the email was sent successfully, false otherwise
	*/
	SendModErrorEmail func(mod_num int, err_str string) bool
}
//////////////////////////////////////////////////////

// _TEMP_DIR is the full path to the main directory of the temporary files.
const _TEMP_DIR string = PersonalConsts.VISOR_DIR + "temp/"
// _DATA_DIR is the full path to the main directory of the data files.
const _DATA_DIR string = PersonalConsts.VISOR_DIR + "data/"
// _DATA_DIR is the full path to the main directory of the modules.
const _MOD_DIR string = PersonalConsts.VISOR_DIR + "Modules/"

// _MOD_FOLDER_PREFFIX is the preffix of the modules' folders.
const _MOD_FOLDER_PREFFIX string = "MOD_"

// _MOD_MAIN_INFO_JSON is the name of the file containing the main module information.
const _MOD_MAIN_INFO_JSON string = "mod_main_info.json"
// _MOD_INFO_JSON is the name of the file containing custom module information.
const _MOD_INFO_JSON string = "mod_info.json"

// _MOD_NUMS_NAMES is a map of the numbers of the modules and their names. Use with the NUM_MOD_ constants.
var _MOD_NUMS_NAMES map[int]string = map[int]string{
	1: "MOD_1",
	2: "S.M.A.R.T. Checker",
	3: "MOD_3",
	4: "RSS Feed Notifier",
	5: "Email Sender",
	6: "MOD_6",
	7: "MOD_7",
}
const NUM_MOD_1 int = 1
const NUM_MOD_SMARTChecker int = 2
const NUM_MOD_3 int = 3
const NUM_MOD_RssFeedNotifier int = 4
const NUM_MOD_EmailSender int = 5
const NUM_MOD_6 int = 6
const NUM_MOD_7 int = 7

// MAX_WAIT_NEXT_TIMESTAMP is the maximum number of seconds to wait for the next timestamp to be registered by a module.
const MAX_WAIT_NEXT_TIMESTAMP int64 = 5

// _RunFileInfo is the struct of the file containing information about the running of a module.
type _RunFileInfo struct {
	// Last_pid is the PID of the last process that ran the module.
	Last_pid          int
	// Last_timestamp_s is the last timestamp in seconds registered by the module.
	Last_timestamp_s int64
}

// ModFileMainInfo is the struct of the file containing the main information about the module.
type ModFileMainInfo struct {
	// Mod_num is the number of the module.
	Mod_num int
	// Run_info is the information about the running of the module.
	Run_info _RunFileInfo
}

// ModInfo is the struct that is provided to a module containing information about it.
type ModInfo struct {
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

> Params:
  - realMain_param_1 – the ModInfo struct of the module
  - realMain_param_2 – the ModFileMainInfo struct of the module
*/
type RealMain func(realMain_param_1 ModInfo, realMain_param_2 ModFileMainInfo)

/*
ModStartup does the startup routine for a module and executes its realMain() function, catching any fatal errors and
sending an email with them.

Call this as the ONLY thing in the main() function of a module.

This function returns when/if the module's realMain() function returns.

Also, this function is outside of the functions struct because it MUST be called before using any of the structs - they
are all initialized by this function.

-----------------------------------------------------------

> Params:
  - mod_num – the number of the module
  - realMain – a pointer to the realMain() function of the module
*/
func ModStartup(mod_num int, realMain RealMain) {
	// Try to run the module, catching any fatal errors and sending an email with them.
	Tcf.Tcf {
		Try: func() {
			// Initialize all the utilities
			InitializeUtils()

			// Module startup routine
			var mod_name string = getModNameMODULES(mod_num)
			printModStartupMODULES(mod_name)
			var modFileInfo ModFileMainInfo = processModRunningMODULES(mod_num)

			// Execute realMain()
			realMain(ModInfo{
				Name:     mod_name,
				Main_Dir: getModDirMODULES(mod_num),
				Data_dir: getModDataDirMODULES(mod_num),
				Temp_dir: getModTempDirMODULES(mod_num),
			}, modFileInfo)
		},
		Catch: func(e Tcf.Exception) {
			var str_email string = ""
			var str_terminal string = ""

			var err, ok = e.(error)
			if ok {
				// tracerr only works with errors
				str_email = getFullErrorMsgGENERAL(err)
				// Colors for the terminal (not for the email because the colors use ANSI escape codes that are read by
				// the terminal only).
				str_terminal = tracerr.SprintSourceColor(tracerr.Wrap(err), 0)
			} else {
				// If the exception is not an error, get general information about it
				var err_str string =
					"Invalid type of error information (not a Go \"error\"). " + getVariableInfoGENERAL(e)
				str_email = err_str
				str_terminal = err_str
			}

			// Send error email and print it
			if !sendModErrorEmailMODULES(mod_num, str_email) {
				fmt.Println("Error sending error email")
			}
			fmt.Println(str_terminal)
		},
	}.Do()
}

func getModNameMODULES(mod_num int) string {
	var mod_name, ok = _MOD_NUMS_NAMES[mod_num]
	if !ok {
		return "INVALID MODULE NUMBER"
	}

	return mod_name
}

func sendModErrorEmailMODULES(mod_num int, err_str string) bool {
	var html_message string = "<p>Error occurred on " + getDateTimeStrTIMEDATE() + ":</p>\n<pre>" + err_str + "</pre>"

	var html string = *getModelFileEMAIL(MODEL_FILE_INFO)
	html = strings.ReplaceAll(html, "|3234_HTML_MESSAGE|", html_message)

	return sendEmailEMAIL(prepareEmlEMAIL(EmailInfo{
		Sender:  "VISOR - Info",
		Mail_to: PersonalConsts.MY_EMAIL_ADDR,
		Subject: "Module error on: " + getModNameMODULES(mod_num),
		Html:    html,
	}), PersonalConsts.MY_EMAIL_ADDR)
}

/*
LoopSleep sleeps for the given number of seconds (with a caveat) and updates the ModFileMainInfo file.

If the number of seconds exceeds MAX_WAIT_NEXT_TIMESTAMP, uses the latter is used instead.

-----------------------------------------------------------

> Params:
  - s – the number of seconds to sleep

> Returns:
  - true if the sleep was successful, false otherwise
*/
func (modFileInfo ModFileMainInfo) LoopSleep(s int64) {
	modFileInfo.Run_info.Last_timestamp_s = time.Now().Unix()
	modFileInfo.Update()

	var seconds = s
	if s > MAX_WAIT_NEXT_TIMESTAMP {
		seconds = MAX_WAIT_NEXT_TIMESTAMP
	}
	time.Sleep(time.Duration(seconds) * time.Second)
}

/*
GetModFileInfo gets the information about the module from its info file.

-----------------------------------------------------------

> Params:
  - v – a pointer to the variable where the information will be stored, with the struct in which the file is written in

> Returns:
  - true if the file was read successfully, false otherwise
*/
func (modInfo ModInfo) GetModFileInfo(v any) bool {
	var p_json_file *string = modInfo.Data_dir.Add(_MOD_INFO_JSON).ReadFile()
	if p_json_file == nil {
		return false
	}

	return fromJsonGENERAL([]byte(*p_json_file), v)
}

/*
Update updates the information about the module in its info file.

-----------------------------------------------------------

> Params:
  - mod_num – the number of the module

> Returns:
  - true if the update was successful, false otherwise
 */
func (modFileInfo ModFileMainInfo) Update() bool {
	var json_str string = *toJsonGENERAL(&modFileInfo)

	return getModDataDirMODULES(modFileInfo.Mod_num).Add(_MOD_MAIN_INFO_JSON).WriteTextFile(json_str)
}

/*
getModDirMODULES gets the full path to the directory of a module.

-----------------------------------------------------------

> Params:
  - mod_num – the number of the module

> Returns:
  - the full path to the directory of the module
*/
func getModDirMODULES(mod_num int) GPath {
	return pathFILESDIRS(_MOD_DIR, _MOD_FOLDER_PREFFIX + strconv.Itoa(mod_num) + "/")
}

/*
getModDataDirMODULES gets the full path to the private data directory of a module.

-----------------------------------------------------------

> Params:
  - mod_num – the number of the module

> Returns:
  - the full path to the private data directory of the module
*/
func getModDataDirMODULES(mod_num int) GPath {
	return pathFILESDIRS(_DATA_DIR, _MOD_FOLDER_PREFFIX + strconv.Itoa(mod_num) + "/")
}

/*
getModTempDirMODULES gets the full path to the private temporary directory of a module.

-----------------------------------------------------------

> Params:
  - mod_num – the number of the module

> Returns:
  - the full path to the private temporary directory of the module
*/
func getModTempDirMODULES(mod_num int) GPath {
	return pathFILESDIRS(_TEMP_DIR, _MOD_FOLDER_PREFFIX + strconv.Itoa(mod_num) + "/")
}

/*
printModStartupMODULES prints the startup message of a module.

-----------------------------------------------------------

> Params:
  - mod_name – the name of the module
*/
func printModStartupMODULES(mod_name string) {
	fmt.Println("--- " + mod_name + " ---")
	fmt.Println("V.I.S.O.R. Systems")
	fmt.Println("---")
	fmt.Println()
}

/*
processModRunningMODULES checks if the module is already running and exits if it is, and if it's not, writes the
necessary information to the module info file.

-----------------------------------------------------------

> Params:
  - mod_num – the number of the module

> Returns:
  - the information about the module
*/
func processModRunningMODULES(mod_num int) ModFileMainInfo {
	var p_modFileInfo *ModFileMainInfo = getModFileInfoMODULES(mod_num)
	if p_modFileInfo == nil {
		goto writeFile
	}

	// Check PID and timestamp
	if isPidRunningPROCESSES(p_modFileInfo.Run_info.Last_pid) &&
			(time.Now().Unix() - p_modFileInfo.Run_info.Last_timestamp_s) < MAX_WAIT_NEXT_TIMESTAMP {

		// todo This is temporary, to see if the modules are being started many times in a row
		sendModErrorEmailMODULES(mod_num, "Module already running")

		fmt.Println("Already running. Exiting...")
		os.Exit(0)
	}

	writeFile:
	var modFileInfo ModFileMainInfo = ModFileMainInfo{
		Mod_num: mod_num,
		Run_info: _RunFileInfo{
			Last_pid:         os.Getpid(),
			Last_timestamp_s: time.Now().Unix(),
		},
	}
	modFileInfo.Update()

	return modFileInfo
}

/*
getModFileInfoMODULES gets the information of a module from the module info file.

-----------------------------------------------------------

> Params:
  - mod_num – the number of the module

> Returns:
  - the information of the module or nil if the module info file doesn't exist or is invalid
 */
func getModFileInfoMODULES(mod_num int) *ModFileMainInfo {
	var p_info *string = getModDataDirMODULES(mod_num).Add(_MOD_MAIN_INFO_JSON).ReadFile()
	if nil == p_info {
		return nil
	}

	var info ModFileMainInfo = ModFileMainInfo{}
	if fromJsonGENERAL([]byte(*p_info), &info) {
		return &info
	}

	return nil
}
