package Utils

/*
InitializeUtils initializes the Utils package.

Call this function at the start of the program, BEFORE using any of the struct Utils functions.
 */
func InitializeUtils() {
	UEmail = _Email_s{
		GetModelFile:           getModelFileEMAIL,
		QueueEmail:             queueEmailEMAIL,
		SendEmail:              sendEmailEMAIL,
		ToQuotedPrintableEMAIL: ToQuotedPrintableEMAIL,
	}

	UGeneral = _General_s{
		RandString:      randStringGENERAL,
		FindAllIndexes:  findAllIndexesGENERAL,
		GetFullErrorMsg: getFullErrorMsgGENERAL,
		PanicGENERAL:    panicGENERAL,
		ToJson:          toJsonGENERAL,
		FromJson:        fromJsonGENERAL,
	}

	UModules = _Modules_s{
		GetModName:        getModNameMODULES,
		SendModErrorEmail: sendModErrorEmailMODULES,
	}

	UFilesDirs = _FilesDirs_s{
		Path: pathFILESDIRS,
	}

	UProcesses = _Processes_s{
		IsPidRunning: isPidRunningPROCESSES,
	}

	UShell = _Shell_s{
		ExecCmd: execCmdSHELL,
	}

	USlices = slices_s{
		DelElem:   delElemSLICES,
		AddElem:   addElemSLICES,
		CopyOuter: copyOuterSLICES,
		CopyFull:  copyFullSLICES,
	}

	UDateTime = _DateTime_s{
		GetDateTimeStr: getDateTimeStrTIMEDATE,
		GetTimeStr:     getTimeStrTIMEDATE,
		GetDateStr:     getDateStrTIMEDATE,
	}

	UWebpages = _Webpages_s{
		GetPageHtml: getPageHtmlTIMEDATE,
	}
}
